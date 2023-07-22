package register

import (
	"onlinedisk-user-service/handler"

	pb "github.com/coderc/onlinedisk-util/grpc/user_service_proto"
	"google.golang.org/grpc"
)

func Register(server *grpc.Server) {
	pb.RegisterUserServiceServer(server, &handler.UserServiceServerImpl{})
}
