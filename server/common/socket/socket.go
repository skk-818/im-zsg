package socket

import (
	"errors"
	"io"
	"net"
	"server/common/libnet"
	"strings"
	"time"
)

// Server 服务器 tcp 服务
type Server struct {
	Name         string          // 服务器名称
	Manager      *libnet.Manager // 连接管理器
	Listener     net.Listener    // 监听器
	Protocol     libnet.Protocol // 协议
	SendChanSize int             // 发送通道大小
}

// NewServer 创建服务器
func NewServer(name string, l net.Listener, p libnet.Protocol, sendChanSize int) *Server {
	return &Server{
		Name:         name,
		Manager:      libnet.NewManager(name),
		Listener:     l,
		Protocol:     p,
		SendChanSize: sendChanSize,
	}
}

// Accept 接收连接
func (s *Server) Accept() (*libnet.Session, error) {
	var tempDelay time.Duration

	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			var ne net.Error
			if errors.As(err, &ne) && ne.Timeout() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if maxDelay := 1 * time.Second; tempDelay > maxDelay {
					tempDelay = maxDelay
				}
				time.Sleep(tempDelay)
				continue
			}
			if strings.Contains(err.Error(), "use of closed network connection") {
				return nil, io.EOF
			}
			return nil, err
		}
		// 创建会话
		return libnet.NewSession(s.Manager, s.Protocol.NewCodec(conn), s.SendChanSize), nil
	}
}

// Close 关闭服务器
func (s *Server) Close() {
	_ = s.Listener.Close()
	s.Manager.Close()
}

// NewServe 创建服务器
func NewServe(name, address string, protocol libnet.Protocol, sendChanSize int) (*Server, error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}
	return NewServer(name, listener, protocol, sendChanSize), nil
}
