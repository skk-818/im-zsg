package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	RpcClient    zrpc.RpcClientConf
	KqConf       kq.KqConf
	KqEtcd       discov.EtcdConf
	Name         string
	TCPListenOn  string
	WSListenOn   string
	SendChanSize int
}
