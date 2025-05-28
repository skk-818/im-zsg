package server

import (
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/net/websocket"
	"net/http"
	"server/app/edge/client"
	"server/app/edge/internal/svc"
	"server/common/socketio"
)

// WsServer websocket server
type WsServer struct {
	Server *socketio.Server
	svcCtx *svc.ServiceContext
}

func NewWsServer(svcCtx *svc.ServiceContext) *WsServer {
	return &WsServer{
		svcCtx: svcCtx,
	}
}

// Start 启动服务
func (ws *WsServer) Start() {
	err := http.ListenAndServe(ws.Server.Address, nil)
	if err != nil {
		panic(err)
	}
}

// HandleRequest 处理请求
func (ws *WsServer) HandleRequest(conn *websocket.Conn) {
	session, err := ws.Server.Accept(conn)
	if err != nil {
		panic(err)
	}
	cli := client.NewClient(ws.Server.Manager, session, ws.svcCtx.RpcClient)
	ws.sessionLoop(cli)
}

// sessionLoop 会话循环
func (ws *WsServer) sessionLoop(client *client.Client) {
	message, err := client.Receive()
	if err != nil {
		logx.Errorf("[ws:sessionLoop] client.Receive error: %v", err)
		_ = client.Close()
		return
	}
	// login check
	err = client.Login(message)
	if err != nil {
		logx.Errorf("[ws:sessionLoop] client.Login error: %v", err)
		_ = client.Close()
		return
	}

	//client.HeartBeat() ws 本身具有心跳机制 不需要 client.HeartBeat()

	for {
		message, err = client.Receive()
		if err != nil {
			logx.Errorf("[ws:sessionLoop] client.Receive error: %v", err)
			_ = client.Close()
			return
		}
		err = client.HandlePackage(message)
		if err != nil {
			logx.Errorf("[ws:sessionLoop] client.HandleMessage error: %v", err)
			_ = client.Close()
			return
		}
	}
}
