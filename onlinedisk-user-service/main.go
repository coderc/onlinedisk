package main

import (
	"fmt"
	"net"
	"onlinedisk-user-service/register"

	"github.com/coderc/onlinedisk-util/config"
	"github.com/coderc/onlinedisk-util/db"
	"github.com/coderc/onlinedisk-util/logger"
	"google.golang.org/grpc"
)

func Init() {
	config.GetConfig().Init("./config/config.yaml")
	logger.GetLogger().Init(&config.GetConfig().LoggerConfig)
	db.GetSingleton().Init(&config.GetConfig().DBConfig)
}

func main() {
	Init()

	// 启动rpc server
	rpcServer := grpc.NewServer()
	// 注册服务
	register.Register(rpcServer)
	lister, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", config.GetConfig().ServerConfig.Port))
	if err != nil {
		panic(err)
	}
	// 监听端口
	err = rpcServer.Serve(lister)
	if err != nil {
		panic(err)
	}
}
