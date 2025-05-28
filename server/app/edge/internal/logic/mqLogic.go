package logic

import (
	"context"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"google.golang.org/protobuf/proto"
	"server/app/edge/internal/svc"
	"server/app/rpc/rpc"
	"server/common/libnet"
	"server/common/session"
	"server/common/socket"
	"server/common/socketio"
)

// MqLogic 实现了 kq.Handler 接口，用于处理 kafka 消息
// 消费 kafka 消息，并将其通过session发送给客户端
type MqLogic struct {
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	server   *socket.Server
	wsServer *socketio.Server
	logx.Logger
}

// NewMqLogic 创建一个 MqLogic 实例
func NewMqLogic(ctx context.Context, svcCtx *svc.ServiceContext, srv *socket.Server, wsSrv *socketio.Server) *MqLogic {
	return &MqLogic{
		ctx:      ctx,
		svcCtx:   svcCtx,
		server:   srv,
		wsServer: wsSrv,
		Logger:   logx.WithContext(ctx),
	}
}

// Consume 消费 kafka 消息
// 消费 kafka 消息，并将其通过session发送给客户端
func (l *MqLogic) Consume(ctx context.Context, _, val string) error {
	// TODO: 消费 kafka 消息，并将其通过session发送给客户端
	logx.Infof("[edge] 消费kafka消息 Consume ==> data %s", val)
	var msg rpc.PostMsg
	if err := proto.Unmarshal([]byte(val), &msg); err != nil {
		logx.Errorf("[Consume] proto.Unmarshal val: %s error: %v", val, err)
		return err
	}
	logx.Infof("[Consume] succ msg: %+v body: %s", msg, msg.Msg)

	if len(msg.ToToken) > 0 {
		sessions := l.server.Manager.GetTokenSessions(msg.ToToken)
		for i := range sessions {
			if sessions[i] == nil {
				logx.Errorf("[Consume] session not found, msg: %v", msg)
				continue
			}
			err := sessions[i].Send(makeMessage(&msg))
			if err != nil {
				logx.Errorf("[Consume] session send error, msg: %v, err: %v", msg, err)
			}
		}
	} else {
		sess := l.server.Manager.GetSession(session.FromString(msg.SessionId))
		wsSess := l.wsServer.Manager.GetSession(session.FromString(msg.SessionId))
		if sess == nil && wsSess == nil {
			logx.Errorf("[Consume] session not found, msg: %v", msg)
			return nil
		}
		if sess != nil {
			err := sess.Send(makeMessage(&msg))
			if err != nil {
				logx.Errorf("[Consume] session send error, msg: %v, err: %v", msg, err)
			}
		}
		if wsSess != nil {
			err := wsSess.Send(makeMessage(&msg))
			if err != nil {
				logx.Errorf("[Consume] session send error, msg: %v, err: %v", msg, err)
			}
		}
	}

	return nil
}

// Consumers 返回一个 kq.Handler 实例
func Consumers(ctx context.Context, svcCtx *svc.ServiceContext, srv *socket.Server, wsSrv *socketio.Server) []service.Service {
	return []service.Service{
		kq.MustNewQueue(svcCtx.Config.KqConf, NewMqLogic(ctx, svcCtx, srv, wsSrv)),
	}
}

// makeMessage 将 rpc 消息转换为 libnet 消息
func makeMessage(msg *rpc.PostMsg) libnet.Message {
	var message libnet.Message
	message.Version = uint8(msg.Version)
	message.Status = uint8(msg.Status)
	message.ServiceId = uint16(msg.ServiceId)
	message.Cmd = uint16(msg.Cmd)
	message.Seq = msg.Seq
	message.Body = []byte(msg.Msg)
	return message
}
