package user_handler

import (
	"context"
	"net/http"

	requestUtil "github.com/coderc/onlinedisk-util/request"

	"github.com/coderc/onlinedisk-util/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	userServiceConn "onlinedisk-backend/grpc/client"

	userServicePb "github.com/coderc/onlinedisk-util/grpc/user_service_proto"

	responseUtil "github.com/coderc/onlinedisk-util/response"
)

func SignupHandler(c *gin.Context) {
	var userInfo requestUtil.RequestUserInfo
	if err := c.BindJSON(&userInfo); err != nil {
		logger.Zap().Warn(errorGetUserInfo, zap.Error(err))
		responseUtil.SendResponse(c, http.StatusBadRequest, responseUtil.SignupFailedCode, nil)
		return
	}

	req := &userServicePb.UserSignUpRequest{
		Username:        userInfo.Username,
		Password:        userInfo.Password,
		ConfirmPassword: userInfo.ConfirmPassword,
	}

	conn := userServiceConn.GetConn()
	userServiceClient := userServicePb.NewUserServiceClient(conn)
	resp, err := userServiceClient.SignUp(context.Background(), req)
	if err != nil {
		logger.Zap().Error(err.Error())
		responseUtil.SendResponse(c, http.StatusInternalServerError, responseUtil.SignupFailedCode, nil)
		return
	}

	if resp.Code != 0 {
		responseUtil.SendResponse(c, http.StatusForbidden, responseUtil.SignupFailedCode, nil)
		return
	}

	responseUtil.SendResponse(c, http.StatusOK, responseUtil.SuccessCode, gin.H{
		"message": "注册成功",
	})
}
