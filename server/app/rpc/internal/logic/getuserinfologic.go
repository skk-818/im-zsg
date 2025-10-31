package logic

import (
	"context"

	"server/app/rpc/internal/svc"
	"server/app/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户相关接口
func (l *GetUserInfoLogic) GetUserInfo(in *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.GetUserInfoResponse{}, nil
}
