package discovery

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

// QueueMaster 队列主控器
// 功能：管理队列配置
type QueueMaster struct {
	members  map[string]kq.KqConf // key 和 kq.KqConf 的映射
	cli      *clientv3.Client     // etcd 客户端
	rootPath string               // key 前缀
	observer QueueObserver        // 观察者
}

// NewQueueMaster 创建队列主控器
// 参数：
// rootPath: key 前缀
// address: etcd 地址
// 返回值：
// *QueueMaster: 队列主控器
// error: 错误
func NewQueueMaster(rootPath string, address []string) (*QueueMaster, error) {
	cfg := clientv3.Config{
		Endpoints:   address,
		DialTimeout: time.Second * 3,
	}
	cli, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}
	return &QueueMaster{
		members:  make(map[string]kq.KqConf),
		cli:      cli,
		rootPath: rootPath,
	}, nil
}

// Register 注册观察者
// 参数：
// o: 观察者
func (m *QueueMaster) Register(o QueueObserver) {
	m.observer = o
}

// notifyUpdate 通知观察者更新
// 参数：
// key: 队列名称
// kqConf: 队列配置
func (m *QueueMaster) notifyUpdate(key string, kqConf kq.KqConf) {
	m.observer.Update(key, kqConf)
}

// notifyDelete 通知观察者删除
// 参数：
// key: 队列名称
func (m *QueueMaster) notifyDelete(key string) {
	m.observer.Delete(key)
}

// addQueueWorker 添加队列
// 参数：
// key: 队列名称
// kqConf: 队列配置
func (m *QueueMaster) addQueueWorker(key string, kqConf kq.KqConf) {
	if len(kqConf.Brokers) == 0 || len(kqConf.Topic) == 0 {
		logx.Errorf("invalid kqConf: %+v", kqConf)
		return
	}
	m.members[key] = kqConf
	m.notifyUpdate(key, kqConf)
}

// updateQueueWorker 更新队列
// 参数：
// key: 队列名称
// kqConf: 队列配置
func (m *QueueMaster) updateQueueWorker(key string, kqConf kq.KqConf) {
	if len(kqConf.Brokers) == 0 || len(kqConf.Topic) == 0 {
		logx.Errorf("invalid kqConf: %+v", kqConf)
		return
	}
	m.members[key] = kqConf
	m.notifyUpdate(key, kqConf)
}

// deleteQueueWorker 删除队列
// 参数：
// key: 队列名称
func (m *QueueMaster) deleteQueueWorker(key string) {
	delete(m.members, key)
	m.notifyDelete(key)
}

// watchQueueWorkers 监听队列
func (m *QueueMaster) watchQueueWorkers() {
	watcher := clientv3.NewWatcher(m.cli)
	defer watcher.Close()

	// 监听队列
	watchChan := watcher.Watch(context.Background(), m.rootPath, clientv3.WithPrefix())
	for resp := range watchChan { // 监听队列
		for _, event := range resp.Events { // 从 watchChan 中获取事件
			switch event.Type {
			case clientv3.EventTypePut:
				// 添加队列
				var kqConf kq.KqConf
				err := json.Unmarshal(event.Kv.Value, &kqConf)
				if err != nil {
					logx.Errorf("unmarshal kqConf failed, err: %v", err)
					continue
				}
				if event.IsCreate() { // 创建 队列
					m.addQueueWorker(string(event.Kv.Key), kqConf)
				} else if event.IsModify() { // 修改 队列
					m.updateQueueWorker(string(event.Kv.Key), kqConf)
				}
			case clientv3.EventTypeDelete:
				// 删除队列
				m.deleteQueueWorker(string(event.Kv.Key))
			}
		}
	}
}
