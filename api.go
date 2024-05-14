package m7s

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	. "github.com/shirou/gopsutil/v3/net"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gopkg.in/yaml.v3"
	"m7s.live/m7s/v5/pb"
	"m7s.live/m7s/v5/pkg"
	"m7s.live/m7s/v5/pkg/config"
	"m7s.live/m7s/v5/pkg/util"
)

var localIP string
var empty = &emptypb.Empty{}

func (s *Server) SysInfo(context.Context, *emptypb.Empty) (res *pb.SysInfoResponse, err error) {
	if localIP == "" {
		if conn, err := net.Dial("udp", "114.114.114.114:80"); err == nil {
			localIP, _, _ = strings.Cut(conn.LocalAddr().String(), ":")
		}
	}
	res = &pb.SysInfoResponse{
		Version:   Version,
		LocalIP:   localIP,
		StartTime: timestamppb.New(s.StartTime),
	}
	return
}

func (s *Server) StreamInfo(ctx context.Context, req *pb.StreamSnapRequest) (res *pb.StreamInfoResponse, err error) {
	// s.Call(func() {
	// 	if pub, ok := s.Streams.Get(req.StreamPath); ok {
	// 		res = &pb.StreamInfoResponse{
	// 		}
	// 	} else {
	// 		err = pkg.ErrNotFound
	// 	}
	// })
	return
}

func (s *Server) AudioTrackSnap(ctx context.Context, req *pb.StreamSnapRequest) (res *pb.AudioTrackSnapShotResponse, err error) {
	// s.Call(func() {
	// 	if pub, ok := s.Streams.Get(req.StreamPath); ok {
	// 		res = pub.AudioSnapShot()
	// 	} else {
	// 		err = pkg.ErrNotFound
	// 	}
	// })
	return
}

func (s *Server) VideoTrackSnap(ctx context.Context, req *pb.StreamSnapRequest) (res *pb.VideoTrackSnapShotResponse, err error) {
	s.Call(func() {
		if pub, ok := s.Streams.Get(req.StreamPath); ok {
			res = &pb.VideoTrackSnapShotResponse{}
			if !pub.VideoTrack.IsEmpty() {
				vcc := pub.VideoTrack.AVTrack.ICodecCtx.(pkg.IVideoCodecCtx)
				res.Width = uint32(vcc.GetWidth())
				res.Height = uint32(vcc.GetHeight())
				res.Info = pub.VideoTrack.GetInfo()
				pub.VideoTrack.Ring.Next().Do(func(v *pkg.AVFrame) {
					var snap pb.TrackSnapShot
					snap.Sequence = v.Sequence
					snap.Timestamp = uint32(v.Timestamp / time.Millisecond)
					snap.WriteTime = timestamppb.New(v.WriteTime)
					snap.Wrap = make([]*pb.Wrap, len(v.Wraps))
					snap.KeyFrame = v.IDR
					for i, wrap := range v.Wraps {
						snap.Wrap[i] = &pb.Wrap{
							Timestamp: uint32(wrap.GetTimestamp() / time.Millisecond),
							Size:      uint32(wrap.GetSize()),
							Data:      wrap.String(),
						}
					}
					res.Ring = append(res.Ring, &snap)
				})
			}
		} else {
			err = pkg.ErrNotFound
		}
	})
	return
}

func (s *Server) Restart(ctx context.Context, req *pb.RequestWithId) (res *emptypb.Empty, err error) {
	if Servers[req.Id] != nil {
		Servers[req.Id].Stop(pkg.ErrRestart)
	}
	return empty, err
}

func (s *Server) Shutdown(ctx context.Context, req *pb.RequestWithId) (res *emptypb.Empty, err error) {
	if Servers[req.Id] != nil {
		Servers[req.Id].Stop(pkg.ErrStopFromAPI)
	} else {
		return nil, pkg.ErrNotFound
	}
	return empty, err
}

func (s *Server) StopSubscribe(ctx context.Context, req *pb.StopSubscribeRequest) (res *pb.StopSubscribeResponse, err error) {
	s.Call(func() {
		if subscriber, ok := s.Subscribers.Get(int(req.Id)); ok {
			subscriber.Stop(errors.New("stop by api"))
		} else {
			err = pkg.ErrNotFound
		}
	})
	return &pb.StopSubscribeResponse{
		Success: err == nil,
	}, err
}

// /api/stream/list
func (s *Server) StreamList(_ context.Context, req *pb.StreamListRequest) (res *pb.StreamListResponse, err error) {
	s.Call(func() {
		var streams []*pb.StreamSummay
		for _, publisher := range s.Streams.Items {
			var audioTrack, videoTrack string
			var bps int32
			if !publisher.VideoTrack.IsEmpty() {
				bps += int32(publisher.VideoTrack.AVTrack.BPS)
				videoTrack = publisher.VideoTrack.FourCC().String()
			}
			if !publisher.AudioTrack.IsEmpty() {
				bps += int32(publisher.AudioTrack.AVTrack.BPS)
				audioTrack = publisher.AudioTrack.FourCC().String()
			}
			streams = append(streams, &pb.StreamSummay{
				Path:        publisher.StreamPath,
				State:       int32(publisher.State),
				StartTime:   timestamppb.New(publisher.StartTime),
				Subscribers: int32(len(publisher.Subscribers)),
				AudioTrack:  audioTrack,
				VideoTrack:  videoTrack,
				Bps:         bps,
				Type:        publisher.Plugin.Meta.Name,
			})
		}
		res = &pb.StreamListResponse{List: streams, Total: int32(s.Streams.Length), PageNum: req.PageNum, PageSize: req.PageSize}
	})
	return
}

func (s *Server) API_Summary_SSE(rw http.ResponseWriter, r *http.Request) {
	util.ReturnFetchValue(func() *pb.SummaryResponse {
		ret, _ := s.Summary(r.Context(), nil)
		return ret
	}, rw, r)
}

func (s *Server) Summary(context.Context, *emptypb.Empty) (res *pb.SummaryResponse, err error) {
	s.Call(func() {
		dur := time.Since(s.lastSummaryTime)
		if dur < time.Second {
			res = s.lastSummary
			return
		}
		v, _ := mem.VirtualMemory()
		d, _ := disk.Usage("/")
		nv, _ := IOCounters(true)
		res = &pb.SummaryResponse{
			Memory: &pb.Usage{
				Total: v.Total >> 20,
				Free:  v.Available >> 20,
				Used:  v.Used >> 20,
				Usage: float32(v.UsedPercent),
			},
			HardDisk: &pb.Usage{
				Total: d.Total >> 30,
				Free:  d.Free >> 30,
				Used:  d.Used >> 30,
				Usage: float32(d.UsedPercent),
			},
		}
		if cc, _ := cpu.Percent(time.Second, false); len(cc) > 0 {
			res.CpuUsage = float32(cc[0])
		}
		netWorks := []*pb.NetWorkInfo{}
		for i, n := range nv {
			info := &pb.NetWorkInfo{
				Name:    n.Name,
				Receive: n.BytesRecv,
				Sent:    n.BytesSent,
			}
			if s.lastSummary != nil && len(s.lastSummary.NetWork) > i {
				info.ReceiveSpeed = (n.BytesRecv - s.lastSummary.NetWork[i].Receive) / uint64(dur.Seconds())
				info.SentSpeed = (n.BytesSent - s.lastSummary.NetWork[i].Sent) / uint64(dur.Seconds())
			}
			netWorks = append(netWorks, info)
		}
		res.StreamCount = int32(s.Streams.Length)
		res.NetWork = netWorks
		s.lastSummary = res
		s.lastSummaryTime = time.Now()
	})
	return
}

// /api/config/json/{name}
func (s *Server) api_Config_JSON_(rw http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	var conf *config.Config
	if name == "global" {
		conf = &s.Config
	} else {
		p, ok := s.Plugins.Get(name)
		if !ok {
			http.Error(rw, pkg.ErrNotFound.Error(), http.StatusNotFound)
			return
		}
		conf = &p.Config
	}
	rw.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(rw).Encode(conf.GetMap())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) GetConfig(_ context.Context, req *pb.GetConfigRequest) (res *pb.GetConfigResponse, err error) {
	res = &pb.GetConfigResponse{}
	var conf *config.Config
	if req.Name == "global" {
		conf = &s.Config
	} else {
		p, ok := s.Plugins.Get(req.Name)
		if !ok {
			err = pkg.ErrNotFound
			return
		}
		conf = &p.Config
	}
	var mm []byte
	mm, err = yaml.Marshal(conf.File)
	if err != nil {
		return
	}
	res.File = string(mm)

	mm, err = yaml.Marshal(conf.Modify)
	if err != nil {
		return
	}
	res.Modified = string(mm)

	mm, err = yaml.Marshal(conf.GetMap())
	if err != nil {
		return
	}
	res.Merged = string(mm)
	return
}

func (s *Server) ModifyConfig(_ context.Context, req *pb.ModifyConfigRequest) (res *pb.ModifyConfigResponse, err error) {
	var conf *config.Config
	if req.Name == "global" {
		conf = &s.Config
		defer s.SaveConfig()
	} else {
		p, ok := s.Plugins.Get(req.Name)
		if !ok {
			err = pkg.ErrNotFound
			return
		}
		defer p.SaveConfig()
		conf = &p.Config
	}
	var modified map[string]any
	err = yaml.Unmarshal([]byte(req.Yaml), &modified)
	if err != nil {
		return
	}
	conf.ParseModifyFile(modified)
	return
}
