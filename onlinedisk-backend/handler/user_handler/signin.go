package user_handler

import (
	"context"
	"net/http"
	userServiceConn "onlinedisk-backend/grpc/client"

	requestUtil "github.com/coderc/onlinedisk-util/request"
	responseUtil "github.com/coderc/onlinedisk-util/response"

	"github.com/coderc/onlinedisk-util/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	userServicePb "github.com/coderc/onlinedisk-util/grpc/user_service_proto"
)

const (
	errorGetUserInfo             = "get user info failed"
	errorUsernameInvalid         = "username invalid"
	errorPasswordInvalidInSignup = "password invalid"
	errorCreateUUIDFailed        = "create uuid failed"
	errorInsertUserInfoFailed    = "insert user info failed"
)

func SigninHandler(c *gin.Context) {
	var userInfo requestUtil.RequestUserInfo
	if err := c.BindJSON(&userInfo); err != nil {
		logger.Zap().Warn(errorGetUserInfo, zap.Error(err))
		responseUtil.SendResponse(c, http.StatusBadRequest, responseUtil.SigninFailedCode, nil)
		return
	}

	req := &userServicePb.UserSignInRequest{
		Username: userInfo.Username,
		Password: userInfo.Password,
	}

	conn := userServiceConn.GetConn()
	userServiceClient := userServicePb.NewUserServiceClient(conn)
	resp, err := userServiceClient.SignIn(context.TODO(), req)
	if err != nil {
		logger.Zap().Error(err.Error())
		responseUtil.SendResponse(c, http.StatusInternalServerError, responseUtil.SigninFailedCode, nil)
		return
	}

	if resp.Code != 0 {
		responseUtil.SendResponse(c, http.StatusForbidden, responseUtil.SigninFailedCode, nil)
		return

	}

	// 返回token
	logger.Zap().Debug("signin success", zap.String("username", userInfo.Username))
	responseUtil.SendResponse(c, http.StatusOK, responseUtil.SuccessCode, gin.H{
		"token": resp.Token,
		"userModel": gin.H{
			"id":          resp.Id,
			"uuid":        resp.Uuid,
			"username":    userInfo.Username,
			"create_time": resp.CreateTime,
			"update_time": resp.UpdateTime,
		},
	})
}
