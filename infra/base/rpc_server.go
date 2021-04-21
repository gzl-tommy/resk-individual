package base

import (
	"github.com/sirupsen/logrus"
	"net"
	"net/rpc"
	"reflect"
	"gzl-tommy/resk-individual/infra"
)

var rpcServer *rpc.Server

func RpcServer() *rpc.Server {
	Check(rpcServer)
	return rpcServer
}

func RpcRegister(ri interface{}) {
	logrus.Infof("goRPC Register:%s", reflect.TypeOf(ri).String())
	RpcServer().Register(ri)
}

type GoRPCStarter struct {
	infra.BaseStarter
	server *rpc.Server
}

func (s *GoRPCStarter) Init(ctx infra.StarterContext) {
	s.server = rpc.NewServer()
	rpcServer = s.server
}

func (s *GoRPCStarter) Start(ctx infra.StarterContext) {
	port := ctx.Props().GetDefault("app.rpc.port", "18082")
	// 监听网络断开
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Infof("tcp port listened for rpc:%s", port)

	// 处理网络连接和请求
	go s.server.Accept(listener)
}
