package file_handler

import (
	"net/http"
	"onlinedisk-backend/pkg/file_store"
	"strconv"

	"github.com/coderc/onlinedisk-util/logger"
	"github.com/gin-gonic/gin"

	responseUtil "github.com/coderc/onlinedisk-util/response"
)

// FileDownloadHandler 用户端下载文件
func FileDownloadHandler(c *gin.Context) {
	fileName := c.Request.URL.Query().Get("filename")
	fileUUIDStr := c.Request.URL.Query().Get("uuid")
	fileUUID, err := strconv.ParseInt(fileUUIDStr, 10, 64)
	if err != nil {
		logger.Zap().Error(err.Error())
		responseUtil.SendResponse(c, http.StatusBadRequest, responseUtil.FileDownloadFailedCode, nil)
		return
	}

	fileBytes, err := file_store.Download(fileUUID)
	if err != nil {
		logger.Zap().Error(err.Error())
		responseUtil.SendResponse(c, http.StatusInternalServerError, responseUtil.FileDownloadFailedCode, nil)
		return
	}

	// 设置响应头
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Data(http.StatusOK, "application/octet-stream", fileBytes)
}
