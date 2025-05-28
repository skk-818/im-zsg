package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/collection"
	"google.golang.org/protobuf/proto"
	"server/app/rpc/internal/svc"
	"server/app/rpc/pb"
	"server/common/session"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPostMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostMessageLogic {
	return &PostMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// PostMessage implements the RPC method PostMessage
// 功能：接收消息 将消息投递到 kafka 中
func (l *PostMessageLogic) PostMessage(in *pb.PostMsg) (*pb.PostResponse, error) {

	logx.Infof("[rpc] 接收到消息 ==> PostMessage: %v", in)

	var (
		allDevice bool
		name      string
		token     string
		id        uint64
	)

	// 这里的 token 是接收方的 token
	if len(in.ToToken) != 0 {
		token = in.Token
		allDevice = true
	} else {
		name, token, id = session.FromString(in.SessionId).Info()
	}

	// 取接受方的 sessionId 列表,这里的 session 是接收方的 用户 session
	sessionIds, err := l.svcCtx.BizRedis.Zrange(token, 0, -1)
	if err != nil {
		return nil, err
	}
	if len(sessionIds) == 0 {
		return nil, err
	}

	// 修改一下消息内容
	in.Msg = "feedback: " + in.Msg
	data, err := proto.Marshal(in)
	if err != nil {
		return nil, err
	}

	// 遍历 sessionId 列表，根据 session
	set := collection.NewSet() // 用来存放接收方的 edge 名字

	for _, sessionId := range sessionIds {
		respName, _, respId := session.FromString(sessionId).Info()
		if allDevice {
			set.Add(respName)
		} else {
			if name == respName && id == respId { // 单设备
				if queueWorker, ok := l.svcCtx.QueueList.Load(respName); ok {
					err = queueWorker.Push(context.Background(), string(data))
					if err != nil {
						logx.Errorf("[PostMessage] push data: %s error: %v", string(data), err)
						return nil, err
					}
				} else {
					logx.Severe("invalid session")
				}
			} else {
				logx.Severe("invalid session")
			}
		}
	}

	if set.Count() > 0 { // 给所有设备发消息
		logx.Infof("send to %d devices", set.Count())
	}

	// 循环 set
	for _, respName := range set.KeysStr() { // 遍历 set 中的元素
		if queueWorker, ok := l.svcCtx.QueueList.Load(respName); ok {
			err = queueWorker.Push(context.Background(), string(data))
			if err != nil {
				return nil, err
			}
		} else {
			logx.Errorf("invalid session")
		}
	}

	return &pb.PostResponse{}, nil
}
