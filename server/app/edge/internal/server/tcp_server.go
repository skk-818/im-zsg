package server

import (
	"github.com/zeromicro/go-zero/core/logx"
	"server/app/edge/client"
	"server/app/edge/internal/svc"
	"server/common/discovery"
	"server/common/socket"
)

// TcpServer  tcp server
type TcpServer struct {
	svcCtx *svc.ServiceContext
	Server *socket.Server
}

func NewTcpServer(svcCtx *svc.ServiceContext) *TcpServer {
	return &TcpServer{
		svcCtx: svcCtx,
	}
}

// HandleRequest 处理请求
func (srv *TcpServer) HandleRequest() {
	for {
		session, err := srv.Server.Accept()
		if err != nil {
			panic(err)
		}
		cli := client.NewClient(srv.Server.Manager, session, srv.svcCtx.RpcClient)
		go srv.sessionLoop(cli)
	}
}

// sessionLoop 针对这个客户端对象进行消息处理
func (srv *TcpServer) sessionLoop(cli *client.Client) {
	message, err := cli.Receive()
	if err != nil {
		logx.Errorf("[sessionLoop] client.Receive error: %v", err)
		_ = cli.Close()
		return
	}

	// 登录校验
	err = cli.Login(message)
	if err != nil {
		logx.Errorf("[sessionLoop] client.Login error: %v", err)
		_ = cli.Close()
		return
	}

	// 保持心跳
	go cli.HeartBeat()

	// 处理消息
	for {
		message, err := cli.Receive()
		if err != nil {
			logx.Errorf("[sessionLoop] client.Receive error: %v", err)
			_ = cli.Close()
			return
		}
		err = cli.HandlePackage(message)
		if err != nil {
			logx.Errorf("[sessionLoop] client.HandleMessage error: %v", err)
		}
	}
}

// KqHeart 注册kq心跳
func (srv *TcpServer) KqHeart() {
	worker := discovery.NewQueueWorker(srv.svcCtx.Config.KqEtcd.Key, srv.svcCtx.Config.KqEtcd.Hosts, srv.svcCtx.Config.KqConf)
	worker.HeartBeat()
}
