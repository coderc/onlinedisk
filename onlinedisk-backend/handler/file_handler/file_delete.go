package file_handler

import (
	"net/http"
	"onlinedisk-backend/pkg/file_store"
	resp "onlinedisk-backend/response_builder"
	"strconv"

	"github.com/coderc/onlinedisk-util/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	uploadPath               = "./static/file/"
	errorReadFileFromRequest = "read file from request failed"
	errorOpenFile            = "open file failed"
	errorCreateFile          = "create file failed"
	successUploadFile        = "upload file success"
	successUploadFileChunk   = "upload file chunk success"
	errorSelectFile          = "select file failed"
	successDeleteFile        = "delete file success"
)

// FileDeleteHandler 用户端删除文件
func FileDeleteHandler(c *gin.Context) {
	fileUUIDStr := c.GetHeader("uuid")
	fileUUID, err := strconv.ParseInt(fileUUIDStr, 10, 64)
	userUUIDBytes, _ := c.Get("uuid")
	userUUID := userUUIDBytes.(int64)

	if err != nil {
		logger.Zap().Error(err.Error())
		resp.SendResponse(c, http.StatusBadRequest, resp.FileDeleteFailedCode, nil)
		return
	}

	err = file_store.Delete(userUUID, fileUUID)
	if err != nil {
		logger.Zap().Error(err.Error())
		resp.SendResponse(c, http.StatusInternalServerError, resp.FileDeleteFailedCode, nil)
		return
	}

	logger.Zap().Info(successDeleteFile, zap.Int64("fileUUID", fileUUID))
	resp.SendResponse(c, http.StatusOK, resp.SuccessCode, nil)
}
