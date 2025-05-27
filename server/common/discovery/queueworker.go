package discovery

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

// QueueWorker is a queue worker.
// 功能：将 kafka 的配置注册到 etcd 中， 并且实现续租
type QueueWorker struct {
	key    string           // 将这个 kafka 的配置注册到 etcd 中的 key
	kqConf kq.KqConf        // kafka 配置
	client *clientv3.Client // etcd 客户端
}

// NewQueueWorker creates a new queue worker.
func NewQueueWorker(key string, endpoints []string, kqConf kq.KqConf) *QueueWorker {
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second * 3,
	}
	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		panic(err)
	}
	return &QueueWorker{
		key:    key,
		client: etcdClient,
		kqConf: kqConf,
	}
}

// HeartBeat is a heartbeat.
func (q *QueueWorker) HeartBeat() {
	value, err := json.Marshal(q.kqConf)
	if err != nil {
		panic(err)
	}
	q.register(string(value))
}

// register is a register.
func (q *QueueWorker) register(value string) {
	//申请一个45秒的租约
	leaseGrantResp, err := q.client.Grant(context.TODO(), 45)
	if err != nil {
		panic(err)
	}
	//拿到租约的id
	leaseId := leaseGrantResp.ID
	logx.Infof("查看leaseId:%x", leaseId)

	//获得kv api子集
	kv := clientv3.NewKV(q.client)

	//put一个kv，让它与租约关联起来，从而实现10秒后自动过期
	putResp, err := kv.Put(context.TODO(), q.key, value, clientv3.WithLease(leaseId))
	if err != nil {
		panic(err)
	}
	logx.Infof("putResp:%v", putResp)

	//续租，让当前的kv关联的租约不断续租，10秒一次
	keepRespChan, err := q.client.KeepAlive(context.TODO(), leaseId)
	if err != nil {
		panic(err)
	}

	//处理续租应答的协程
	go func() {
		for {
			select {
			case keepResp, ok := <-keepRespChan:
				if !ok {
					logx.Infof("租约已经失效:%x", leaseId)
					q.register(value) // 续租的执行函数
					return
				} else { //每秒会续租一次，所以就会受到一次应答
					logx.Infof("收到自动续租应答:%x", keepResp.ID)
				}
			}
		}
	}()
	// 写入成功
	logx.Info("写入成功:", putResp.Header.Revision)
}
