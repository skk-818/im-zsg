package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/threading"
	"server/app/rpc/internal/config"
	"server/common/discovery"
)

type ServiceContext struct {
	Config    config.Config
	BizRedis  *redis.Redis
	QueueList *QueueList
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 程序执行 第一次初始化
	queueList := GetQueueList(c.QueueEtcd) // 获取队列列表 从 etcd 中获取 并将队列信息保存到本地

	// 开启协程 监听队列
	threading.GoSafe(func() {
		discovery.QueueDiscoveryProc(c.QueueEtcd, queueList) // 实时监听队列，将 etcd 中的队列信息同步到本地
	})
	rds, err := redis.NewRedis(redis.RedisConf{
		Host: c.BizRedis.Host,
		Pass: c.BizRedis.Pass,
		Type: c.BizRedis.Type,
	})
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:    c,
		BizRedis:  rds,
		QueueList: queueList,
	}
}
