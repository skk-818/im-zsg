package logic

import (
	"context"
	"time"

	"server/app/rpc/internal/svc"
	"server/app/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Login 登录
// 功能：登录 将用户信息存入Redis sessionId
// 参数：
// in: 登录请求
// 返回值：
// *pb.LoginResponse 登录响应
// *error 错误
func (l *LoginLogic) Login(in *pb.LoginRequest) (*pb.LoginResponse, error) {
	// todo: add your logic here and delete this line

	_, err := l.svcCtx.BizRedis.Zadd(in.Token, time.Now().UnixMilli(), in.SessionId)
	if err != nil {
		logx.Errorf("[Login] Zadd token: %s sessionId: %s  error: %v", in.Token, in.SessionId, err)
		return nil, err
	}
	_ = l.svcCtx.BizRedis.Expire(in.Token, 3600)

	return &pb.LoginResponse{}, nil
}
