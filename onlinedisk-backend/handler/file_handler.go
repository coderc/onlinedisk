package handler

import (
	"io"
	"net/http"
	resp "onlinedisk-backend/response_builder"
	"os"

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
)

// UploadFileHandler 用户端上传文件
func UploadFileHandler(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		logger.GetLogger().Error(errorReadFileFromRequest, zap.Error(err))
		resp.BuildResponse(c, http.StatusBadRequest, resp.FileUploadFailedCode, nil)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		logger.GetLogger().Error(errorOpenFile, zap.Error(err))
		resp.BuildResponse(c, http.StatusBadRequest, resp.FileUploadFailedCode, nil)
		return
	}

	defer file.Close()

	filePath := uploadPath + fileHeader.Filename
	newFile, err := os.Create(filePath)
	if err != nil {
		logger.GetLogger().Error(errorCreateFile, zap.Error(err))
		resp.BuildResponse(c, http.StatusBadRequest, resp.FileUploadFailedCode, nil)
		return
	}

	defer newFile.Close()

	io.Copy(newFile, file)

	logger.GetLogger().Info(successUploadFile, zap.String("file", filePath))
	resp.BuildResponse(c, http.StatusOK, resp.SuccessCode, nil)
}

// DownloadFileHandler 用户端下载文件
func DownloadFileHandler(c *gin.Context) {

}

// DeleteFileHandler 用户端删除文件
func DeleteFileHandler(c *gin.Context) {

}

// RenameFileHandler 用户端重命名文件
func RenameFileHandler(c *gin.Context) {

}
