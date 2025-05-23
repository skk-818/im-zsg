package logic

import (
	"context"

	"server/app/rpc/internal/svc"
	"server/app/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// @desc ping
func (l *PingLogic) Ping(in *pb.Nil) (*pb.Nil, error) {
	// todo: add your logic here and delete this line

	return &pb.Nil{}, nil
}
