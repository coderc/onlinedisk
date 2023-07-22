package handler

import (
	"context"
	"fmt"

	pb "github.com/coderc/onlinedisk-util/grpc/user_service_proto"
	"github.com/coderc/onlinedisk-util/logger"
	"github.com/coderc/onlinedisk-util/mapper"
	responseUtil "github.com/coderc/onlinedisk-util/response"
	"github.com/coderc/onlinedisk-util/snowflake"
	"github.com/coderc/onlinedisk-util/utils"
	"go.uber.org/zap"
)

type UserServiceServerImpl struct {
	pb.UserServiceServer
}

func (s *UserServiceServerImpl) SignUp(ctx context.Context, req *pb.UserSignUpRequest) (*pb.UserSignUpResponse, error) {
	// 检查用户名是否合法
	if ok := checkUsername(req.Username); !ok {
		logger.Zap().Error(responseUtil.ErrorUsernameInvalid)
		return nil, fmt.Errorf(responseUtil.ErrorUsernameInvalid)
	}

	// 检查密码是否合法
	if ok := checkPasswordInSignup(req.Password, req.ConfirmPassword); !ok {
		logger.Zap().Error(responseUtil.ErrorPasswordInvalid)
		return nil, fmt.Errorf(responseUtil.ErrorPasswordInvalid)
	}

	// 对密码进行加密
	req.Password = utils.EncryptStrMD5(req.Password)

	// 生成 uuid
	uuid, err := snowflake.GetId(1, 1)
	if err != nil {
		logger.Zap().Error(responseUtil.ErrorCreateUUIDFailed, zap.Error(err))
		return nil, fmt.Errorf(responseUtil.ErrorCreateUUIDFailed)
	}

	// 保存用户信息
	err = mapper.InsertUser(uuid, req.Username, req.Password)
	if err != nil {
		logger.Zap().Error(responseUtil.ErrorInsertUserInfoFailed, zap.Error(err))
		return nil, fmt.Errorf(responseUtil.ErrorInsertUserInfoFailed)
	}

	return &pb.UserSignUpResponse{
		Code: responseUtil.SccessCode,
	}, nil
}
