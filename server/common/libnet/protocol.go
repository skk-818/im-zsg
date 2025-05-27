package libnet

import (
	"fmt"
	"net"
	"time"
)

// Codec 编解码器
type Codec interface {
	SetReadDeadline(t time.Time) error  // 设置读取超时时间
	SetWriteDeadline(t time.Time) error // 设置写入超时时间
	Receive() (*Message, error)         // 接收消息
	Send(Message) error                 // 发送消息
	Close() error                       // 关闭连接
}

// Protocol 协议
type Protocol interface {
	NewCodec(conn net.Conn) Codec
}

// Header 消息头
type Header struct {
	Version   uint8
	Status    uint8
	ServiceId uint16
	Cmd       uint16
	Seq       uint32
}

// Message 消息结构体
type Message struct {
	Header
	Body []byte
}

// Format 格式化消息
func (m *Message) Format() string {
	return fmt.Sprintf("Version:%d, Status:%d, ServiceId:%d, Cmd:%d, Seq:%d, Body:%s",
		m.Version, m.Status, m.ServiceId, m.Cmd, m.Seq, string(m.Body))
}
