package register

import (
	"google.golang.org/grpc"
	pb "onlinedisk-user-service/grpc/user_service_proto"
	"onlinedisk-user-service/handler"
)

func Register(server *grpc.Server) {
	pb.RegisterUserServiceServer(server, &handler.UserServiceServerImpl{})
}
