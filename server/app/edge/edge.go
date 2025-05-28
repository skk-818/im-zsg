package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/net/websocket"
	"net/http"
	"server/app/edge/internal/logic"
	"server/app/edge/internal/server"
	"server/common/libnet"
	"server/common/socket"
	"server/common/socketio"

	"server/app/edge/internal/config"
	"server/app/edge/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	zeroservice "github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/edge.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	srvCtx := svc.NewServiceContext(c)

	logx.DisableStat()

	// 长连接
	tcpServer := server.NewTcpServer(srvCtx)
	// websocket
	wsServer := server.NewWsServer(srvCtx)

	protocol := libnet.NewIMProtocol()

	tcpSocket, err := socket.NewServe(srvCtx.Config.Name, srvCtx.Config.TCPListenOn, protocol, srvCtx.Config.SendChanSize)
	if err != nil {
		panic(err)
	}
	tcpServer.Server = tcpSocket

	wsSocket, err := socketio.NewServe(srvCtx.Config.Name, srvCtx.Config.WSListenOn, protocol, srvCtx.Config.SendChanSize)
	if err != nil {
		panic(err)
	}
	wsServer.Server = wsSocket

	// websocket 处理
	http.Handle("/ws", websocket.Handler(func(conn *websocket.Conn) {
		conn.PayloadType = websocket.BinaryFrame
		wsServer.HandleRequest(conn)
	}))

	go wsServer.Start()          // 启动websocket服务
	go tcpServer.HandleRequest() // 启动tcp服务 监听 处理连接请求
	go tcpServer.KqHeart()       // 注册kq心跳

	// 启动服务
	fmt.Printf("Starting tcp server at %s, ws server at: %s...\n", c.TCPListenOn, c.WSListenOn)

	serviceGroup := zeroservice.NewServiceGroup()
	defer serviceGroup.Stop()

	// 启动 kq 消费
	for _, mq := range logic.Consumers(context.Background(), srvCtx, tcpServer.Server, wsServer.Server) {
		serviceGroup.Add(mq)
	}
	serviceGroup.Start()
}
