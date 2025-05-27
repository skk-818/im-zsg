package config

import (
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	BizRedis struct {
		Host string
		Pass string
		Type string
	}
	QueueEtcd discov.EtcdConf
}
