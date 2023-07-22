package handler

import (
	"context"
	"github.com/coderc/onlinedisk-util/jwt"
	"github.com/coderc/onlinedisk-util/logger"
	"github.com/coderc/onlinedisk-util/mapper"
	responseUtil "github.com/coderc/onlinedisk-util/response"
	"github.com/coderc/onlinedisk-util/utils"
	"go.uber.org/zap"
	pb "onlinedisk-user-service/grpc/user_service_proto"
	"time"
)

func (s *UserServiceServerImpl) SignIn(ctx context.Context, req *pb.UserSignInRequest) (*pb.UserSignInResponse, error) {
	// 获取加密后的密码
	req.Password = utils.EncryptStrMD5(req.Password)

	// 获取用户信息
	userModel, err := mapper.QueryUser(req.Username, req.Password)
	if err != nil {
		logger.Zap().Warn(err.Error(), zap.String("username", req.Username))
		return nil, err
	}

	// 生成token
	token, err := jwt.CreateToken(userModel.UUID, time.Now().Add(24*60*60*time.Second).Unix())
	if err != nil {
		logger.Zap().Error(err.Error(), zap.String("username", req.Username))
		return nil, err
	}

	// 返回token
	logger.Zap().Debug("signin success", zap.String("username", req.Username))
	return &pb.UserSignInResponse{
		Code:  responseUtil.SccessCode,
		Token: token,
	}, nil
}
