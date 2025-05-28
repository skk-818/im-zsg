package client

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
	"server/app/rpc/rpc"
	"server/common/libnet"
	"time"
)

// Client 客户端连接对象
type Client struct {
	Session   *libnet.Session
	Manager   *libnet.Manager
	IMRpc     rpc.Rpc
	heartbeat chan *libnet.Message
}

// NewClient  创建客户端连接对象
func NewClient(manager *libnet.Manager, session *libnet.Session, rpc rpc.Rpc) *Client {
	return &Client{
		Session:   session,
		Manager:   manager,
		IMRpc:     rpc,
		heartbeat: make(chan *libnet.Message),
	}
}

// Send  发送消息
func (c *Client) Send(msg libnet.Message) error {
	return c.Session.Send(msg)
}

// Receive  接收消息
func (c *Client) Receive() (*libnet.Message, error) {
	return c.Session.Receive()
}

// Close 关闭连接
func (c *Client) Close() error {
	return c.Session.Close()
}

// Login 登录
func (c *Client) Login(msg *libnet.Message) error {
	loginReq, err := makeLoginMessage(msg)
	if err != nil {
		return err
	}

	c.Session.SetToken(loginReq.Token)
	c.Manager.AddSession(c.Session)

	_, err = c.IMRpc.Login(context.Background(), &rpc.LoginRequest{
		Token:         loginReq.Token,
		Authorization: loginReq.Authorization,
		SessionId:     c.Session.Session().String(),
	})
	if err != nil {
		msg.Status = 1
		msg.Body = []byte(err.Error())
		e := c.Send(*msg)
		if e != nil {
			logx.Errorf("[Login] client.Send error: %v", e)
		}
		return err
	}
	msg.Status = 0
	msg.Body = []byte("登录成功")
	err = c.Send(*msg)
	if err != nil {
		logx.Errorf("[Login] client.Send error: %v", err)
	}

	return err
}

func makeLoginMessage(msg *libnet.Message) (*rpc.LoginRequest, error) {
	// 登录功能还没做
	// 这里临时处理，先把PostMsg中的Msg转换成LoginRequest中的Token和Authorization用于登录处理
	var postMsg rpc.PostMsg
	err := proto.Unmarshal(msg.Body, &postMsg)
	if err != nil {
		return nil, err
	}

	loginReq := rpc.LoginRequest{
		Token:         postMsg.Token,
		Authorization: postMsg.Token,
	}

	return &loginReq, nil
}

// HandlePackage 处理消息
func (c *Client) HandlePackage(msg *libnet.Message) error {
	// 消息转发
	req := makePostMessage(c.Session.Session().String(), msg)
	if req == nil {
		return nil
	}
	_, err := c.IMRpc.PostMessage(context.Background(), req)
	if err != nil {
		logx.Errorf("[HandlePackage] client.PostMessage error: %v", err)
	}
	return err
}

// makePostMessage 将 libnet.Message 转换成 rpc.PostMsg
func makePostMessage(sessionId string, msg *libnet.Message) *rpc.PostMsg {
	var postMessageReq rpc.PostMsg
	err := proto.Unmarshal(msg.Body, &postMessageReq)
	if err != nil {
		logx.Errorf("[makePostMessage] proto.Unmarshal msg: %v error: %v", msg, err)
		return nil
	}
	postMessageReq.Version = uint32(msg.Version)
	postMessageReq.Status = uint32(msg.Status)
	postMessageReq.ServiceId = uint32(msg.ServiceId)
	postMessageReq.Cmd = uint32(msg.Cmd)
	postMessageReq.Seq = msg.Seq
	postMessageReq.SessionId = sessionId

	return &postMessageReq
}

const heartBeatTimeout = time.Second * 60

// HeartBeat 维持心跳
func (c *Client) HeartBeat() error {
	timer := time.NewTimer(heartBeatTimeout)
	for {
		select {
		case heartbeat := <-c.heartbeat:
			_ = c.Session.SetReadDeadline(time.Now().Add(heartBeatTimeout * 5))
			_ = c.Send(*heartbeat)
			break
		case <-timer.C:
		}
	}
}
