package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"server/app/edge/internal/config"
	"server/app/rpc/pb"
	"server/app/rpc/rpc"
)

type ServiceContext struct {
	Config    config.Config
	RpcClient pb.RpcClient
}

func NewServiceContext(c config.Config) *ServiceContext {

	client := zrpc.MustNewClient(c.RpcClient)

	return &ServiceContext{
		Config:    c,
		RpcClient: rpc.NewRpc(client),
	}
}
