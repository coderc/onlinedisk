package file_handler

import (
	"net/http"

	"github.com/coderc/onlinedisk-util/logger"
	"github.com/coderc/onlinedisk-util/mapper"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	responseUtil "github.com/coderc/onlinedisk-util/response"
)

func FileListHandler(c *gin.Context) {
	userUUIDBytes, _ := c.Get("uuid")
	userUUID := userUUIDBytes.(int64)

	userFileModels, err := mapper.QueryUserFileByUserId(userUUID)
	if err != nil {
		logger.Zap().Error(err.Error())
		responseUtil.SendResponse(c, http.StatusInternalServerError, responseUtil.FailedSelectFileListCode, nil)
		return
	}

	logger.Zap().Info("select file list success",
		zap.Int64("userUUID", userUUID), zap.Int("fileCount", len(userFileModels)))
	responseUtil.SendResponse(c, http.StatusOK, responseUtil.SuccessCode, userFileModels)
}
