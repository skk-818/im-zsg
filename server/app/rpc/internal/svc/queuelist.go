package svc

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
	"time"
)

// QueueList 队列列表
// 功能：维护队列列表
// 参数：
// kqs: 队列列表
// l: 锁
type QueueList struct {
	kqs map[string]*kq.Pusher // kqs 队列列表
	l   sync.Mutex            // l 锁
}

// Update 更新队列
// 功能：更新队列
// 参数：
// s: 队列名称
// conf: 队列配置
// 返回值：
// 无
func (q *QueueList) Update(s string, conf kq.KqConf) {
	// 对应 edge 的队列
	pusher := kq.NewPusher(conf.Brokers, conf.Topic)
	q.l.Lock()
	q.kqs[s] = pusher
	q.l.Unlock()
}

// Delete 删除队列
// 功能：删除队列
// 参数：
// s: 队列名称
// 返回值：
// 无
func (q *QueueList) Delete(s string) {
	q.l.Lock()
	delete(q.kqs, s)
	q.l.Unlock()
}

// NewQueueList 创建队列列表
// 功能：创建队列列表
// 参数：
// 无
// 返回值：
// *QueueList 队列列表
func NewQueueList() *QueueList {
	return &QueueList{
		kqs: make(map[string]*kq.Pusher),
	}
}

// GetQueueList 第一次获取队列列表
// 功能：第一次获取队列列表，根据配置文件中的队列配置，创建队列列表并存在内存中
// 参数：
// conf: 配置
// 返回值：
// *QueueList 队列列表
func GetQueueList(conf discov.EtcdConf) *QueueList {

	// 创建队列列表
	queueList := NewQueueList()

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   conf.Hosts,
		DialTimeout: time.Second * 3,
	})
	if err != nil {
		panic(err)
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()

	// 根据前缀获取队列配置
	res, err := cli.Get(ctx, conf.Key, clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}

	// 遍历队列配置
	for _, kv := range res.Kvs {
		// 解析队列配置
		var kqConf kq.KqConf
		err := json.Unmarshal(kv.Value, &kqConf)
		if err != nil {
			logx.Errorf(" invalid data key is: %s value is: %s", string(kv.Key), string(kv.Value))
			continue
		}

		// 判断队列配置是否合法
		if len(kqConf.Brokers) == 0 || len(kqConf.Topic) == 0 {
			continue
		}

		// 创建队列
		edgeQueue := kq.NewPusher(kqConf.Brokers, kqConf.Topic)

		// 将队列添加到队列列表中
		queueList.l.Lock()
		queueList.kqs[string(kv.Key)] = edgeQueue
		queueList.l.Unlock()

	}

	return queueList
}

// Load 加载队列
// 功能：加载队列
// 参数：
// s: 队列名称
// 返回值：
// *kq.Pusher 队列
func (q *QueueList) Load(s string) (*kq.Pusher, bool) {
	q.l.Lock()
	defer q.l.Unlock()

	queue, ok := q.kqs[s]
	return queue, ok
}
