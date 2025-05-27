package discovery

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/discov"
)

// QueueObserver 观察者接口
// 功能：更新队列配置
// 参数：
// key: 队列名称
// kqConf: 队列配置
type QueueObserver interface {
	Update(string, kq.KqConf)
	Delete(string)
}

// QueueDiscoveryProc 队列发现
// 功能：启动队列发现
// 参数：
// conf: 配置
// qo: 观察者
func QueueDiscoveryProc(conf discov.EtcdConf, qo QueueObserver) {
	master, err := NewQueueMaster(conf.Key, conf.Hosts)
	if err != nil {
		panic(err)
	}
	master.Register(qo)        // 注册观察者
	master.watchQueueWorkers() // 监听队列
}
