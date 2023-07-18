package file_handler

import (
	"net/http"

	resp "onlinedisk-backend/response_builder"

	"github.com/coderc/onlinedisk-util/logger"
	"github.com/coderc/onlinedisk-util/mapper"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func FileListHandler(c *gin.Context) {
	userUUIDBytes, _ := c.Get("uuid")
	userUUID := userUUIDBytes.(int64)

	userFileModels, err := mapper.QueryUserFileByUserId(userUUID)
	if err != nil {
		logger.Zap().Error(err.Error())
		resp.SendResponse(c, http.StatusInternalServerError, resp.FailedSelectFileListCode, nil)
		return
	}

	logger.Zap().Info("select file list success",
		zap.Int64("userUUID", userUUID), zap.Int("fileCount", len(userFileModels)))
	resp.SendResponse(c, http.StatusOK, resp.SuccessCode, userFileModels)
}
