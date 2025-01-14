package rtmp

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"

	"m7s.live/m7s/v5"
)

type Client struct {
	NetStream
	ServerInfo map[string]any
}

func NewPushHandler() m7s.PushHandler {
	return &Client{}
}

func NewPullHandler() m7s.PullHandler {
	return &Client{}
}

func (client *Client) Connect(p *m7s.Client) (err error) {
	chunkSize := 4096
	addr := p.RemoteURL
	u, err := url.Parse(addr)
	if err != nil {
		return err
	}
	ps := strings.Split(u.Path, "/")
	if len(ps) < 3 {
		return errors.New("illegal rtmp url")
	}
	isRtmps := u.Scheme == "rtmps"
	if strings.Count(u.Host, ":") == 0 {
		if isRtmps {
			u.Host += ":443"
		} else {
			u.Host += ":1935"
		}
	}
	var conn net.Conn
	if isRtmps {
		var tlsconn *tls.Conn
		tlsconn, err = tls.Dial("tcp", u.Host, &tls.Config{})
		conn = tlsconn
	} else {
		conn, err = net.Dial("tcp", u.Host)
	}
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			conn.Close()
		}
	}()
	client.NetConnection = NewNetConnection(conn, p.Logger)
	if err = client.ClientHandshake(); err != nil {
		return err
	}
	client.AppName = strings.Join(ps[1:len(ps)-1], "/")
	err = client.SendMessage(RTMP_MSG_CHUNK_SIZE, Uint32Message(chunkSize))
	if err != nil {
		return
	}
	client.WriteChunkSize = chunkSize
	path := u.Path
	if len(u.Query()) != 0 {
		path += "?" + u.RawQuery
	}
	err = client.SendMessage(RTMP_MSG_AMF0_COMMAND, &CallMessage{
		CommandMessage{"connect", 1},
		map[string]any{
			"app":      client.AppName,
			"flashVer": "monibuca/" + m7s.Version,
			"swfUrl":   addr,
			"tcUrl":    strings.TrimSuffix(addr, path) + "/" + client.AppName,
		},
		nil,
	})
	for err != nil {
		msg, err := client.RecvMessage()
		if err != nil {
			return err
		}
		switch msg.MessageTypeID {
		case RTMP_MSG_AMF0_COMMAND:
			cmd := msg.MsgData.(Commander).GetCommand()
			switch cmd.CommandName {
			case "_result":
				client.ServerInfo = msg.MsgData.(*ResponseMessage).Properties
				response := msg.MsgData.(*ResponseMessage)
				if response.Infomation["code"] == NetConnection_Connect_Success {

				} else {
					return err
				}
			default:
				fmt.Println(cmd.CommandName)
			}
		}
	}
	client.Info("connect", "remoteURL", p.RemoteURL)
	return
}

func (puller *Client) Pull(p *m7s.Puller) (err error) {
	p.MetaData = puller.ServerInfo
	defer func() {
		puller.Close()
		if p := recover(); p != nil {
			err = p.(error)
		}
		p.Dispose(err)
	}()
	err = puller.SendMessage(RTMP_MSG_AMF0_COMMAND, &CommandMessage{"createStream", 2})
	for err == nil {
		msg, err := puller.RecvMessage()
		if err != nil {
			return err
		}
		switch msg.MessageTypeID {
		case RTMP_MSG_AUDIO:
			p.WriteAudio(msg.AVData.WrapAudio())
		case RTMP_MSG_VIDEO:
			p.WriteVideo(msg.AVData.WrapVideo())
		case RTMP_MSG_AMF0_COMMAND:
			cmd := msg.MsgData.(Commander).GetCommand()
			switch cmd.CommandName {
			case "_result":
				if response, ok := msg.MsgData.(*ResponseCreateStreamMessage); ok {
					puller.StreamID = response.StreamId
					m := &PlayMessage{}
					m.StreamId = response.StreamId
					m.TransactionId = 4
					m.CommandMessage.CommandName = "play"
					URL, _ := url.Parse(p.Client.RemoteURL)
					ps := strings.Split(URL.Path, "/")
					p.Args = URL.Query()
					m.StreamName = ps[len(ps)-1]
					if len(p.Args) > 0 {
						m.StreamName += "?" + p.Args.Encode()
					}
					puller.SendMessage(RTMP_MSG_AMF0_COMMAND, m)
					// if response, ok := msg.MsgData.(*ResponsePlayMessage); ok {
					// 	if response.Object["code"] == "NetStream.Play.Start" {

					// 	} else if response.Object["level"] == Level_Error {
					// 		return errors.New(response.Object["code"].(string))
					// 	}
					// } else {
					// 	return errors.New("pull faild")
					// }
				}
			}
		}
	}
	return
}

func (pusher *Client) Push(p *m7s.Pusher) (err error) {
	p.MetaData = pusher.ServerInfo
	pusher.SendMessage(RTMP_MSG_AMF0_COMMAND, &CommandMessage{"createStream", 2})
	for {
		msg, err := pusher.RecvMessage()
		if err != nil {
			return err
		}
		switch msg.MessageTypeID {
		case RTMP_MSG_AMF0_COMMAND:
			cmd := msg.MsgData.(Commander).GetCommand()
			switch cmd.CommandName {
			case Response_Result, Response_OnStatus:
				if response, ok := msg.MsgData.(*ResponseCreateStreamMessage); ok {
					pusher.StreamID = response.StreamId
					URL, _ := url.Parse(p.Client.RemoteURL)
					_, streamPath, _ := strings.Cut(URL.Path, "/")
					_, streamPath, _ = strings.Cut(streamPath, "/")
					p.Args = URL.Query()
					if len(p.Args) > 0 {
						streamPath += "?" + p.Args.Encode()
					}
					pusher.SendMessage(RTMP_MSG_AMF0_COMMAND, &PublishMessage{
						CURDStreamMessage{
							CommandMessage{
								"publish",
								1,
							},
							response.StreamId,
						},
						streamPath,
						"live",
					})
				} else if response, ok := msg.MsgData.(*ResponsePublishMessage); ok {
					if response.Infomation["code"] == NetStream_Publish_Start {
						audio, video := pusher.CreateSender(true)
						go m7s.PlayBlock(&p.Subscriber, audio.HandleAudio, video.HandleVideo)
					} else {
						return errors.New(response.Infomation["code"].(string))
					}
				}
			}
		}
	}
}
