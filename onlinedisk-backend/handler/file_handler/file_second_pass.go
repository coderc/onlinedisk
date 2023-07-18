package file_handler

import (
	"net/http"
	"net/url"
	"onlinedisk-backend/pkg/file_store"
	resp "onlinedisk-backend/response_builder"

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
	err := file_store.Chunk(sha1, fileName, userUUID)
	if err != nil { // 秒传失败
		logger.Zap().Error(err.Error())
		resp.SendResponse(c, http.StatusInternalServerError, resp.FileUploadFailedCode, nil)
		return
	}
	// 秒传成功
	logger.Zap().Info(successUploadFileChunk)
	resp.SendResponse(c, http.StatusOK, resp.SuccessCode, nil)
	return
}
