package file_handler

import (
	"net/http"
	"net/url"
	"onlinedisk-backend/pkg/file_store"

	responseUtil "github.com/coderc/onlinedisk-util/response"

	"github.com/coderc/onlinedisk-util/logger"
	"github.com/gin-gonic/gin"
)

// FileSecondPassHandler 秒传接口
func FileSecondPassHandler(c *gin.Context) {
	userUUIDBytes, _ := c.Get("uuid")
	userUUID := userUUIDBytes.(int64)
	sha1 := c.GetHeader("sha1")
	fileName := c.GetHeader("filename")
	fileName, _ = url.QueryUnescape(fileName)
	err := file_store.CheckSecondPass(sha1, fileName, userUUID)
	if err != nil { // 秒传失败
		logger.Zap().Error(err.Error())
		responseUtil.SendResponse(c, http.StatusInternalServerError, responseUtil.FileUploadFailedCode, nil)
		return
	}
	// 秒传成功
	logger.Zap().Info(successUploadFileChunk)
	responseUtil.SendResponse(c, http.StatusOK, responseUtil.SuccessCode, nil)
	return
}
