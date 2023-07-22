package client

import (
	"fmt"
	"sync"

	"github.com/coderc/onlinedisk-util/config"
	"google.golang.org/grpc"
)

var (
	once sync.Once
	cli  *grpc.ClientConn
)

func Init(config *config.UserServiceConfigStruct) error {
	var err error
	cli, err = grpc.Dial(fmt.Sprintf("%s:%d", config.Host, config.Port), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return err
}

func GetConn() *grpc.ClientConn {
	if cli == nil {
		panic("grpc client is nil")
	}
	return cli
}
