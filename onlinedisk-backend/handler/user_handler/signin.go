package user_handler

import (
	"net/http"
	"time"

	requestUtil "github.com/coderc/onlinedisk-util/request"
	responseUtil "github.com/coderc/onlinedisk-util/response"

	"github.com/coderc/onlinedisk-util/mapper"

	"github.com/coderc/onlinedisk-util/jwt"
	"github.com/coderc/onlinedisk-util/logger"
	"github.com/coderc/onlinedisk-util/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

	// 获取加密后的密码
	userInfo.Password = utils.EncryptStrMD5(userInfo.Password)

	// 获取用户信息
	userModel, err := mapper.QueryUser(userInfo.Username, userInfo.Password)
	if err != nil {
		logger.Zap().Warn(err.Error(), zap.String("username", userInfo.Username))
		responseUtil.SendResponse(c, http.StatusBadRequest, responseUtil.SigninFailedCode, nil)
		return
	}

	// 生成token
	token, err := jwt.CreateToken(userModel.UUID, time.Now().Add(24*60*60*time.Second).Unix())
	if err != nil {
		logger.Zap().Error(err.Error(), zap.String("username", userInfo.Username))
		responseUtil.SendResponse(c, http.StatusInternalServerError, responseUtil.SigninFailedCode, nil)
		return
	}

	// 返回token
	logger.Zap().Debug("signin success", zap.String("username", userInfo.Username))
	responseUtil.SendResponse(c, http.StatusOK, responseUtil.SuccessCode, gin.H{
		"token":     token,
		"userModel": userModel,
	})
}
