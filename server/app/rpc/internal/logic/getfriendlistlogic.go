package logic

import (
	"context"

	"server/app/rpc/internal/svc"
	"server/app/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendListLogic {
	return &GetFriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 好友相关接口
func (l *GetFriendListLogic) GetFriendList(in *pb.GetFriendListRequest) (*pb.GetFriendListResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.GetFriendListResponse{}, nil
}
