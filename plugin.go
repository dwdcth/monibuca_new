package m7s

import (
	"context"
	myip "github.com/husanpao/ip"
	"gorm.io/gorm"
	"log/slog"
	"m7s.live/m7s/v5/pkg/db"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	gatewayRuntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
	. "m7s.live/m7s/v5/pkg"
	"m7s.live/m7s/v5/pkg/config"
	"m7s.live/m7s/v5/pkg/util"
)

type DefaultYaml string

type PluginMeta struct {
	Name                string
	Version             string //插件版本
	Type                reflect.Type
	defaultYaml         DefaultYaml //默认配置
	ServiceDesc         *grpc.ServiceDesc
	RegisterGRPCHandler func(context.Context, *gatewayRuntime.ServeMux, *grpc.ClientConn) error
}

func (plugin *PluginMeta) Init(s *Server, userConfig map[string]any) {
	instance, ok := reflect.New(plugin.Type).Interface().(IPlugin)
	if !ok {
		panic("plugin must implement IPlugin")
	}
	p := reflect.ValueOf(instance).Elem().FieldByName("Plugin").Addr().Interface().(*Plugin)
	p.handler = instance
	p.Meta = plugin
	p.Server = s
	p.Logger = s.Logger.With("plugin", plugin.Name)
	p.Context, p.CancelCauseFunc = context.WithCancelCause(s.Context)
	upperName := strings.ToUpper(plugin.Name)
	if os.Getenv(upperName+"_ENABLE") == "false" {
		p.Disabled = true
		p.Warn("disabled by env")
		return
	}
	p.Config.Parse(p.GetCommonConf(), upperName)
	p.Config.Parse(instance, upperName)
	for _, fname := range MergeConfigs {
		if name := strings.ToLower(fname); p.Config.Has(name) {
			p.Config.Get(name).ParseGlobal(s.Config.Get(name))
		}
	}
	if plugin.defaultYaml != "" {
		var defaultConf map[string]any
		if err := yaml.Unmarshal([]byte(plugin.defaultYaml), &defaultConf); err != nil {
			p.Error("parsing default config", "error", err)
		} else {
			p.Config.ParseDefaultYaml(defaultConf)
		}
	}
	p.Config.ParseUserFile(userConfig)
	finalConfig, _ := yaml.Marshal(p.Config.GetMap())
	p.Logger.Handler().(*MultiLogHandler).SetLevel(ParseLevel(p.config.LogLevel))
	p.Debug("config", "detail", string(finalConfig))
	if s.DisableAll {
		p.Disabled = true
	}
	if userConfig["enable"] == false {
		p.Disabled = true
	} else if userConfig["enable"] == true {
		p.Disabled = false
	}
	if p.Disabled {
		p.Warn("plugin disabled")
		return
	} else {
		p.assign()
	}
	p.Info("init", "version", plugin.Version)
	var err error
	if p.config.DSN == s.GetCommonConf().DSN {
		p.DB = s.DB
	} else if p.config.DSN != "" {
		if factory, ok := db.Factory[p.config.DBType]; ok {
			s.DB, err = gorm.Open(factory(p.config.DSN), &gorm.Config{})
			if err != nil {
				s.Error("failed to connect database", "error", err, "dsn", s.config.DSN, "type", s.config.DBType)
				p.Disabled = true
				return
			}
		}
	}
	err = instance.OnInit()
	if err != nil {
		p.Error("init", "error", err)
		p.Stop(err)
		return
	}
	if plugin.ServiceDesc != nil && s.grpcServer != nil {
		s.grpcServer.RegisterService(plugin.ServiceDesc, instance)
		if plugin.RegisterGRPCHandler != nil {
			if err = plugin.RegisterGRPCHandler(p.Context, s.config.HTTP.GetGRPCMux(), s.grpcClientConn); err != nil {
				p.Error("init", "error", err)
				p.Stop(err)
			} else {
				p.Info("grpc handler registered")
			}
		}
	}
	s.Plugins.Add(p)
	p.Start()
}

type iPlugin interface {
	nothing()
}

type IPlugin interface {
	OnInit() error
	OnEvent(any)
}

type IRegisterHandler interface {
	RegisterHandler() map[string]http.HandlerFunc
}

type IPullerPlugin interface {
	NewPullHandler() PullHandler
	GetPullableList() []string
}

type IPusherPlugin interface {
	NewPushHandler() PushHandler
}

type IRecorderPlugin interface {
	NewRecordHandler() RecordHandler
}

type ITCPPlugin interface {
	OnTCPConnect(*net.TCPConn)
}

type IUDPPlugin interface {
	OnUDPConnect(*net.UDPConn)
}

var plugins []PluginMeta

func InstallPlugin[C iPlugin](options ...any) error {
	var c *C
	t := reflect.TypeOf(c).Elem()
	meta := PluginMeta{
		Name: strings.TrimSuffix(t.Name(), "Plugin"),
		Type: t,
	}

	_, pluginFilePath, _, _ := runtime.Caller(1)
	configDir := filepath.Dir(pluginFilePath)

	if _, after, found := strings.Cut(configDir, "@"); found {
		meta.Version = after
	} else {
		meta.Version = pluginFilePath
	}
	for _, option := range options {
		switch v := option.(type) {
		case DefaultYaml:
			meta.defaultYaml = v
		case *grpc.ServiceDesc:
			meta.ServiceDesc = v
		case func(context.Context, *gatewayRuntime.ServeMux, *grpc.ClientConn) error:
			meta.RegisterGRPCHandler = v
		}
	}
	plugins = append(plugins, meta)
	return nil
}

type Plugin struct {
	Unit
	Disabled bool
	Meta     *PluginMeta
	config   config.Common
	config.Config
	handler IPlugin
	Server  *Server
	DB      *gorm.DB
}

func (Plugin) nothing() {

}

func (p *Plugin) GetKey() string {
	return p.Meta.Name
}

func (p *Plugin) GetGlobalCommonConf() *config.Common {
	return p.Server.GetCommonConf()
}

func (p *Plugin) GetCommonConf() *config.Common {
	return &p.config
}

func (p *Plugin) GetHandler() IPlugin {
	return p.handler
}

func (p *Plugin) GetPublicIP(netcardIP string) string {
	if p.config.PublicIP != "" {
		return p.config.PublicIP
	}
	if publicIP, ok := Routes[netcardIP]; ok { //根据网卡ip获取对应的公网ip
		return publicIP
	}
	localIp := myip.InternalIPv4()
	if publicIP, ok := Routes[localIp]; ok {
		return publicIP
	}
	return localIp
}

func (p *Plugin) settingPath() string {
	return filepath.Join(p.Server.SettingDir, strings.ToLower(p.Meta.Name)+".yaml")
}

func (p *Plugin) assign() {
	f, err := os.Open(p.settingPath())
	defer f.Close()
	if err == nil {
		var modifyConfig map[string]any
		err = yaml.NewDecoder(f).Decode(&modifyConfig)
		if err != nil {
			panic(err)
		}
		p.Config.ParseModifyFile(modifyConfig)
	}
	var handlerMap map[string]http.HandlerFunc
	if v, ok := p.handler.(IRegisterHandler); ok {
		handlerMap = v.RegisterHandler()
	}
	p.registerHandler(handlerMap)
}

func (p *Plugin) Stop(err error) {
	p.Unit.Stop(err)
	p.config.HTTP.StopListen()
	p.config.TCP.StopListen()
}

func (p *Plugin) Start() {
	httpConf := &p.config.HTTP
	if httpConf.ListenAddrTLS != "" && (httpConf.ListenAddrTLS != p.Server.config.HTTP.ListenAddrTLS) {
		p.Info("https listen at ", "addr", httpConf.ListenAddrTLS)
		go func() {
			p.Stop(httpConf.ListenTLS())
		}()
	}
	if httpConf.ListenAddr != "" && (httpConf.ListenAddr != p.Server.config.HTTP.ListenAddr) {
		p.Info("http listen at ", "addr", httpConf.ListenAddr)
		go func() {
			p.Stop(httpConf.Listen())
		}()
	}
	tcpConf := &p.config.TCP

	tcphandler, ok := p.handler.(ITCPPlugin)
	if !ok {
		tcphandler = p
	}

	if tcpConf.ListenAddr != "" && tcpConf.AutoListen {
		p.Info("listen tcp", "addr", tcpConf.ListenAddr)
		go func() {
			err := tcpConf.Listen(tcphandler.OnTCPConnect)
			if err != nil {
				p.Error("listen tcp", "addr", tcpConf.ListenAddr, "error", err)
				p.Stop(err)
			}
		}()
	}
	if tcpConf.ListenAddrTLS != "" && tcpConf.AutoListen {
		p.Info("listen tcp tls", "addr", tcpConf.ListenAddrTLS)
		go func() {
			err := tcpConf.ListenTLS(tcphandler.OnTCPConnect)
			if err != nil {
				p.Error("listen tcp tls", "addr", tcpConf.ListenAddrTLS, "error", err)
				p.Stop(err)
			}
		}()
	}
	udpConf := &p.config.UDP

	udpHandler, ok := p.handler.(IUDPPlugin)
	if !ok {
		udpHandler = p
	}

	if udpConf.ListenAddr != "" && udpConf.AutoListen {
		p.Info("listen udp", "addr", udpConf.ListenAddr)
		go func() {
			err := udpConf.Listen(udpHandler.OnUDPConnect)
			if err != nil {
				p.Error("listen udp", "addr", udpConf.ListenAddr, "error", err)
				p.Stop(err)
			}
		}()

	}
}

func (p *Plugin) OnInit() error {
	return nil
}

func (p *Plugin) onEvent(event any) {
	switch v := event.(type) {
	case *Publisher:
		if h, ok := p.handler.(interface{ OnPublish(*Publisher) }); ok {
			h.OnPublish(v)
		}
	case *Puller:
		if h, ok := p.handler.(interface{ OnPull(*Puller) }); ok {
			h.OnPull(v)
		}
	}
	p.handler.OnEvent(event)
}

func (p *Plugin) OnEvent(event any) {

}

func (p *Plugin) OnTCPConnect(conn *net.TCPConn) {
	p.handler.OnEvent(conn)
}

func (p *Plugin) OnUDPConnect(conn *net.UDPConn) {
	p.handler.OnEvent(conn)
}

func (p *Plugin) Publish(streamPath string, options ...any) (publisher *Publisher, err error) {
	publisher = &Publisher{Publish: p.config.Publish}
	if p.config.EnableAuth {
		if onAuthPub, ok := p.Server.OnAuthPubs[p.Meta.Name]; ok {
			authPromise := util.NewPromise(publisher)
			onAuthPub(authPromise)
			if _, err = authPromise.Await(); err != nil {
				p.Warn("auth failed", "error", err)
				return
			}
		}
	}
	for _, option := range options {
		switch v := option.(type) {
		case func(*config.Publish):
			v(&publisher.Publish)
		}
	}
	publisher.Init(p, streamPath, &publisher.Publish, options...)
	_, err = p.Server.Call(publisher)
	return
}

func (p *Plugin) Pull(streamPath string, url string, options ...any) (puller *Puller, err error) {
	puller = &Puller{Pull: p.config.Pull}
	puller.Client.Proxy = p.config.Pull.Proxy
	puller.Client.RemoteURL = url
	puller.Client.PubSubBase = &puller.PubSubBase
	puller.Publish = p.config.Publish
	puller.PublishTimeout = 0
	puller.StreamPath = streamPath
	var pullHandler PullHandler
	for _, option := range options {
		switch v := option.(type) {
		case PullHandler:
			pullHandler = v
		}
	}
	puller.Init(p, streamPath, &puller.Publish, options...)
	if _, err = p.Server.Call(puller); err != nil {
		return
	}
	if v, ok := p.handler.(IPullerPlugin); pullHandler == nil && ok {
		pullHandler = v.NewPullHandler()
	}
	if pullHandler != nil {
		err = puller.Start(pullHandler)
	}
	return
}

func (p *Plugin) Record(streamPath string, filePath string, options ...any) (recorder *Recorder, err error) {
	recorder = &Recorder{
		Record: p.config.Record,
	}
	if err = os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return
	}
	recorder.StreamPath = streamPath
	recorder.Subscribe = p.config.Subscribe
	if recorder.File, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666); err != nil {
		return
	}
	defer func() {
		err = recorder.File.Close()
		if info, err := recorder.File.Stat(); err == nil && info.Size() == 0 {
			os.Remove(recorder.File.Name())
		}
	}()
	recorder.Init(p, streamPath, &recorder.Subscribe, options...)
	if _, err = p.Server.Call(recorder); err != nil {
		return
	}
	recorder.Publisher.WaitTrack()
	var recordHandler RecordHandler
	if v, ok := p.handler.(IRecorderPlugin); recordHandler == nil && ok {
		recordHandler = v.NewRecordHandler()
	}
	if recordHandler != nil {
		err = recorder.Start(recordHandler)
	}
	return
}

func (p *Plugin) Subscribe(streamPath string, options ...any) (subscriber *Subscriber, err error) {
	subscriber = &Subscriber{Subscribe: p.config.Subscribe}
	if p.config.EnableAuth {
		if onAuthSub, ok := p.Server.OnAuthSubs[p.Meta.Name]; ok {
			authPromise := util.NewPromise(subscriber)
			onAuthSub(authPromise)
			if _, err = authPromise.Await(); err != nil {
				p.Warn("auth failed", "error", err)
				return
			}
		}
	}
	for _, option := range options {
		switch v := option.(type) {
		case func(*config.Subscribe):
			v(&subscriber.Subscribe)
		}
	}
	subscriber.Init(p, streamPath, &subscriber.Subscribe, options...)
	if subscriber.Subscribe.BufferTime > 0 {
		subscriber.Subscribe.SubMode = SUBMODE_BUFFER
	}
	_, err = p.Server.Call(subscriber)
	subscriber.Publisher.WaitTrack()
	return
}

func (p *Plugin) registerHandler(handlers map[string]http.HandlerFunc) {
	t := reflect.TypeOf(p.handler)
	v := reflect.ValueOf(p.handler)
	// 注册http响应
	for i, j := 0, t.NumMethod(); i < j; i++ {
		name := t.Method(i).Name
		if name == "ServeHTTP" {
			continue
		}
		switch handler := v.Method(i).Interface().(type) {
		case func(http.ResponseWriter, *http.Request):
			patten := strings.ToLower(strings.ReplaceAll(name, "_", "/"))
			p.handle(patten, http.HandlerFunc(handler))
		}
	}
	for patten, handler := range handlers {
		p.handle(patten, handler)
	}
	if rootHandler, ok := p.handler.(http.Handler); ok {
		p.handle("/", rootHandler)
	}
}

func (p *Plugin) Push(streamPath string, url string, options ...any) (pusher *Pusher, err error) {
	pusher = &Pusher{Push: p.config.Push}
	pusher.Client.PubSubBase = &pusher.PubSubBase
	pusher.Client.Proxy = p.config.Push.Proxy
	pusher.Client.RemoteURL = url
	pusher.Subscribe = p.config.Subscribe
	pusher.StreamPath = streamPath
	var pushHandler PushHandler
	for _, option := range options {
		switch v := option.(type) {
		case PushHandler:
			pushHandler = v
		}
	}
	pusher.Init(p, streamPath, &pusher.Subscribe, options...)
	if _, err = p.Server.Call(pusher); err != nil {
		return
	}
	pusher.Publisher.WaitTrack()
	if v, ok := p.handler.(IPusherPlugin); pushHandler == nil && ok {
		pushHandler = v.NewPushHandler()
	}
	if pushHandler != nil {
		err = pusher.Start(pushHandler)
	}
	return
}

func (p *Plugin) logHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		p.Debug("visit", "path", r.URL.String(), "remote", r.RemoteAddr)
		name := strings.ToLower(p.Meta.Name)
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/"+name)
		handler.ServeHTTP(rw, r)
	})
}

func (p *Plugin) handle(pattern string, handler http.Handler) {
	if p == nil {
		return
	}
	last := pattern == "/"
	if !strings.HasPrefix(pattern, "/") {
		pattern = "/" + pattern
	}
	handler = p.logHandler(handler)
	p.config.HTTP.Handle(pattern, handler, last)
	if p.Server != p.handler {
		pattern = "/" + strings.ToLower(p.Meta.Name) + pattern
		p.Debug("http handle added to Server", "pattern", pattern)
		p.Server.config.HTTP.Handle(pattern, handler, last)
	}
	p.Server.apiList = append(p.Server.apiList, pattern)
}

func (p *Plugin) AddLogHandler(handler slog.Handler) {
	p.Server.LogHandler.Add(handler)
}

func (p *Plugin) PostToServer(event any) {
	if p.Server.eventChan == nil {
		panic("eventChan is nil")
	}
	p.Server.PostMessage(event)
}

func (p *Plugin) SaveConfig() (err error) {
	_, err = p.Server.Call(func() error {
		if p.Modify == nil {
			os.Remove(p.settingPath())
			return nil
		}
		file, err := os.OpenFile(p.settingPath(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err == nil {
			defer file.Close()
			err = yaml.NewEncoder(file).Encode(p.Modify)
		}
		if err == nil {
			p.Info("config saved")
		}
		return err
	})
	return
}
