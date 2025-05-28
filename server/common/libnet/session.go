package libnet

import (
	"errors"
	"math/rand"
	"server/common/session"
	"sync"
	"sync/atomic"
	"time"
)

func init() {
	globalSessionId = uint64(rand.New(rand.NewSource(time.Now().Unix())).Uint32()) // 生成一个随机的 session id 全局的
}

var (
	SessionClosedError  = errors.New("session Closed")
	SessionBlockedError = errors.New("session Blocked")

	globalSessionId uint64
)

// Session 是连接对象
// 发送消息，接收消息，关闭连接 都是通过 Session 来完成的
type Session struct {
	id         uint64       // 连接 id
	token      string       // 连接的 token
	codec      Codec        // 编解码
	manager    *Manager     // 这个连接的管理器，用于注销
	sendChan   chan Message // 消息通道
	closeFlag  int32        // 原子标记是否已关闭
	closeChan  chan int     // 关闭通知通道
	closeMutex sync.Mutex   // 关闭锁
}

// NewSession 创建一个连接对象 每一个连接对象都对应一个 session 对象
// 并开启一个协程用于发送消息
func NewSession(manager *Manager, codec Codec, sendChanSize int) *Session {
	s := &Session{
		codec:     codec,
		manager:   manager,
		closeChan: make(chan int),
		id:        atomic.AddUint64(&globalSessionId, 1),
	}
	if sendChanSize > 0 {
		s.sendChan = make(chan Message, sendChanSize)
		go s.sendLoop() // 启动发送消息的协程
	}

	return s
}

// sendLoop 发送消息
func (s *Session) sendLoop() {
	for {
		select {
		case msg := <-s.sendChan: // 从发送消息通道中获取消息
			if err := s.codec.Send(msg); err != nil { // 发送消息
				s.Close()
				return
			}
		case <-s.closeChan: // 关闭通知
			return
		}
	}
}

// Name 返回连接的名字 edge 的名字
func (s *Session) Name() string {
	return s.manager.Name
}

// ID 返回连接的 id
func (s *Session) ID() uint64 {
	return s.id
}

// Token 返回连接的 token
func (s *Session) Token() string {
	return s.token
}

// Close 关闭连接
func (s *Session) Close() error {
	if atomic.CompareAndSwapInt32(&s.closeFlag, 0, 1) {
		err := s.codec.Close()
		close(s.closeChan)
		if s.manager != nil {
			s.manager.RemoveSession(s)
		}
		return err
	}

	return SessionClosedError
}

// Session 返回连接的用户 session
func (s *Session) Session() session.Session {
	return session.NewSession(s.manager.Name, s.token, s.id)
}

// SetToken 设置连接的 token
func (s *Session) SetToken(token string) {
	s.token = token
}

// Receive 接收消息
func (s *Session) Receive() (*Message, error) {
	return s.codec.Receive()
}

// Send 发送消息
func (s *Session) Send(msg Message) error {
	if s.IsClosed() {
		return SessionClosedError
	}
	if s.sendChan == nil {
		return s.codec.Send(msg)
	}
	select {
	case s.sendChan <- msg:
		return nil
	default:
		return SessionBlockedError
	}
}

// IsClosed 判断连接是否已关闭
func (s *Session) IsClosed() bool {
	return atomic.LoadInt32(&s.closeFlag) == 1
}

// SetReadDeadline 设置读取超时时间
func (s *Session) SetReadDeadline(t time.Time) error {
	return s.codec.SetReadDeadline(t)
}

// SetWriteDeadline 设置写入超时时间
func (s *Session) SetWriteDeadline(t time.Time) error {
	return s.codec.SetWriteDeadline(t)
}
