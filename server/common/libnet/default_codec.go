package libnet

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net"
	"time"
)

const maxBodySize = 1 << 12

/*
总长度
----|
header头长度
--|
header头
-|-|--|--|----|------...
版本号|状态码|消息类型|命令|seq|pb body体

header头长度=1字节版本号+1字节状态码+2字节消息类型+2字节命令+4字节seq
总长度=header头+header头长度+pb body体长度

----|--|-|-|--|--|----|body
总长度|header头长度|版本号|状态码|消息类型|命令|seq｜body
总长度=2+1+1+2+2+4+len(body)
header头长度=1+1+2+2+4
*/

const (
	packSize      = 4
	headerSize    = 2
	verSize       = 1
	statusSize    = 1
	serviceIdSize = 2
	cmdSize       = 2
	seqSize       = 4
	rawHeaderSize = verSize + statusSize + serviceIdSize + cmdSize + seqSize
	maxPackSize   = maxBodySize + rawHeaderSize + headerSize + packSize
	// offset
	headerOffset    = 0
	verOffset       = headerOffset + headerSize
	statusOffset    = verOffset + verSize
	serviceIdOffset = statusOffset + statusSize
	cmdOffset       = serviceIdOffset + serviceIdSize
	seqOffset       = cmdOffset + cmdSize
	bodyOffset      = seqOffset + seqSize
)

var (
	ErrRawPackLen   = errors.New("default server codec pack length error")
	ErrRawHeaderLen = errors.New("default server codec header length error")
)

// NewImCodec 创建 imCodec 实例
func NewImCodec(conn net.Conn) *imCodec {
	return &imCodec{
		conn: conn,
	}
}

// IMProtocol im 协议
type IMProtocol struct{}

func NewIMProtocol() Protocol {
	return &IMProtocol{}
}

// imCodec 实现 Codec 接口
type imCodec struct {
	conn net.Conn
}

func (p *IMProtocol) NewCodec(conn net.Conn) Codec {
	return &imCodec{conn: conn}
}

// -------------------------------------------------------------------------
// 实现 Codec 接口
// -------------------------------------------------------------------------

// SetReadDeadline 设置读取超时时间
func (c *imCodec) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

// SetWriteDeadline 设置写入超时时间
func (c *imCodec) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

// Close 关闭连接
func (c *imCodec) Close() error {
	return c.conn.Close()
}

// readPackSize 读取包长度
func (c *imCodec) readPackSize() (uint32, error) {
	return c.readUint32BE()
}

// readUint32BE 读取uint32大端
func (c *imCodec) readUint32BE() (uint32, error) {
	b := make([]byte, packSize)
	_, err := io.ReadFull(c.conn, b)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(b), nil
}

// readPacket 读取包
func (c *imCodec) readPacket(msgSize uint32) ([]byte, error) {
	b := make([]byte, msgSize)
	_, err := io.ReadFull(c.conn, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Receive 接收消息
func (c *imCodec) Receive() (*Message, error) {
	packLen, err := c.readPackSize() // 读取包长度
	if err != nil {
		return nil, err
	}

	if packLen > maxPackSize { // 包长度超过最大长度
		return nil, errors.New(fmt.Sprintf("default server codec pack length error ==> packLen: %d, maxPackSize: %d", packLen, maxPackSize))
	}

	if packLen < rawHeaderSize { // 包长度小于最小长度
		return nil, errors.New(fmt.Sprintf("default server codec header length error ==> packLen: %d, rawHeaderSize: %d", packLen, rawHeaderSize))
	}

	buf, err := c.readPacket(packLen)
	if err != nil {
		return nil, err
	}

	msg := &Message{}
	headerLen := binary.BigEndian.Uint16(buf[headerOffset:verOffset])
	msg.Version = buf[verOffset]
	msg.Status = buf[statusOffset]
	msg.ServiceId = binary.BigEndian.Uint16(buf[serviceIdOffset:cmdOffset])
	msg.Cmd = binary.BigEndian.Uint16(buf[cmdOffset:seqOffset])
	msg.Seq = binary.BigEndian.Uint32(buf[seqOffset:bodyOffset])
	//logx.Infof("msg.Seq:%+v", msg.Seq)

	if headerLen != rawHeaderSize {
		return nil, errors.New(fmt.Sprintf("default server codec header length error ==> headerLen: %d, rawHeaderSize: %d", headerLen, rawHeaderSize))
	}

	if packLen > uint32(headerLen) {
		msg.Body = buf[bodyOffset:packLen]
	}

	logx.Infof("receive msg:%+v", msg)
	return msg, nil
}

// Send 发送消息
func (c *imCodec) Send(msg Message) error {
	packLen := headerSize + rawHeaderSize + len(msg.Body)
	packLenBuf := make([]byte, packSize)
	binary.BigEndian.PutUint32(packLenBuf[:packSize], uint32(packLen))

	buf := make([]byte, packLen)
	// header
	binary.BigEndian.PutUint16(buf[headerOffset:], uint16(rawHeaderSize))
	buf[verOffset] = msg.Version
	buf[statusOffset] = msg.Status
	binary.BigEndian.PutUint16(buf[serviceIdOffset:], msg.ServiceId)
	binary.BigEndian.PutUint16(buf[cmdOffset:], msg.Cmd)
	binary.BigEndian.PutUint32(buf[seqOffset:], msg.Seq)

	// body
	copy(buf[headerSize+rawHeaderSize:], msg.Body)
	allBuf := append(packLenBuf, buf...)
	n, err := c.conn.Write(allBuf)
	if err != nil {
		return err
	}
	if n != len(allBuf) {
		return fmt.Errorf("n:%d, len(buf):%d", n, len(buf))
	}
	return nil
}
