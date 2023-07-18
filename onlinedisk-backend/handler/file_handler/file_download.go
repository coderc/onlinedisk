package file_handler

import (
	"net/http"
	"onlinedisk-backend/pkg/file_store"
	"strconv"

	resp "onlinedisk-backend/response_builder"

	"github.com/coderc/onlinedisk-util/logger"
	"github.com/gin-gonic/gin"
)

// FileDownloadHandler 用户端下载文件
func FileDownloadHandler(c *gin.Context) {
	fileName := c.Request.URL.Query().Get("filename")
	fileUUIDStr := c.Request.URL.Query().Get("uuid")
	fileUUID, err := strconv.ParseInt(fileUUIDStr, 10, 64)
	if err != nil {
		logger.Zap().Error(err.Error())
		resp.SendResponse(c, http.StatusBadRequest, resp.FileDownloadFailedCode, nil)
		return
	}

	fileBytes, err := file_store.Download(fileUUID)
	if err != nil {
		logger.Zap().Error(err.Error())
		resp.SendResponse(c, http.StatusInternalServerError, resp.FileDownloadFailedCode, nil)
		return
	}

	// 设置响应头
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Data(http.StatusOK, "application/octet-stream", fileBytes)
}
