package logic

import (
	"context"
	"github.com/google/uuid"
	"strings"

	"server/app/rpc/internal/svc"
	"server/app/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAuthLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthLoginLogic {
	return &AuthLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AuthLogin 用户 api 登录
// @desc：登录，验证用户信息 并且为用户派发token
// @params：
// @in: *pb.AuthLoginRequest
// @return： (*pb.AuthLoginResponse, error)
func (l *AuthLoginLogic) AuthLogin(in *pb.AuthLoginRequest) (*pb.AuthLoginResponse, error) {
	tokenStr := uuid.NewString()
	token := strings.ReplaceAll(tokenStr, "-", "")

	return &pb.AuthLoginResponse{
		Token:         token,
		Authorization: token + "jwt",
	}, nil
}
