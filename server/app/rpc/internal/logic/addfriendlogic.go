package logic

import (
	"context"

	"server/app/rpc/internal/svc"
	"server/app/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFriendLogic {
	return &AddFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddFriendLogic) AddFriend(in *pb.AddFriendRequest) (*pb.AddFriendResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.AddFriendResponse{}, nil
}
