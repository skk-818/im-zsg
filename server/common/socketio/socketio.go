package socketio

import (
	"net"
	"server/common/libnet"
)

// Server ws server
type Server struct {
	Name         string
	Address      string
	Manager      *libnet.Manager
	Protocol     libnet.Protocol
	SendChanSize int
}

// NewServe 创建一个 ws server
func NewServe(name, address string, protocol libnet.Protocol, sendChanSize int) (*Server, error) {
	return &Server{
		Name:         name,
		Address:      address,
		Manager:      libnet.NewManager(name),
		Protocol:     protocol,
		SendChanSize: sendChanSize,
	}, nil
}

// Accept 接收一个连接
func (s *Server) Accept(conn net.Conn) (*libnet.Session, error) {
	return libnet.NewSession(
		s.Manager,
		s.Protocol.NewCodec(conn),
		s.SendChanSize,
	), nil
}
