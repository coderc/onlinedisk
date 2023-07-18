package user_handler

import (
	"net/http"
	requestInfo "onlinedisk-backend/request_struct"
	resp "onlinedisk-backend/response_builder"

	"github.com/coderc/onlinedisk-util/logger"
	"github.com/coderc/onlinedisk-util/mapper"
	"github.com/coderc/onlinedisk-util/snowflake"
	"github.com/coderc/onlinedisk-util/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SignupHandler(c *gin.Context) {
	var userInfo requestInfo.RequestUserInfo
	if err := c.BindJSON(&userInfo); err != nil {
		logger.Zap().Warn(errorGetUserInfo, zap.Error(err))
		resp.SendResponse(c, http.StatusBadRequest, resp.SignupFailedCode, nil)
		return
	}

	// 判断用户名是否合法
	if ok := checkUsername(userInfo.Username); !ok {
		logger.Zap().Warn(errorUsernameInvalid, zap.String("username", userInfo.Username))
		resp.SendResponse(c, http.StatusBadRequest, resp.SignupFailedCode, nil)
		return
	}

	// 判断密码是否合法
	if ok := checkPasswordInSignup(userInfo.Password, userInfo.ConfirmPassword); !ok {
		logger.Zap().Warn(errorPasswordInvalidInSignup, zap.String("password", userInfo.Password), zap.String("confirmPassword", userInfo.ConfirmPassword))
		resp.SendResponse(c, http.StatusBadRequest, resp.SignupFailedCode, nil)
		return
	}

	// 对密码进行加密
	userInfo.Password = utils.EncryptStrMD5(userInfo.Password)

	// 生成 uuid
	uuid, err := snowflake.GetId(1, 1)
	if err != nil {
		logger.Zap().Error(errorCreateUUIDFailed, zap.Error(err))
		resp.SendResponse(c, http.StatusInternalServerError, resp.SignupFailedCode, nil)
		return
	}

	// 保存用户信息
	err = mapper.InsertUser(uuid, userInfo.Username, userInfo.Password)
	if err != nil {
		logger.Zap().Error(errorInsertUserInfoFailed, zap.Error(err))
		resp.SendResponse(c, http.StatusInternalServerError, resp.SignupFailedCode, nil)
		return
	}
	resp.SendResponse(c, http.StatusOK, resp.SuccessCode, gin.H{
		"message": "注册成功",
	})
}
