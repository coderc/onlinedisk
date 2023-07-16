package handler

import (
	"io"
	"net/http"
	"net/url"
	"onlinedisk-backend/pkg/file_store"
	resp "onlinedisk-backend/response_builder"
	"os"

	"github.com/coderc/onlinedisk-util/logger"
	"github.com/coderc/onlinedisk-util/mapper"
	"github.com/coderc/onlinedisk-util/model"
	"github.com/coderc/onlinedisk-util/snowflake"
	"github.com/coderc/onlinedisk-util/utils"
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
)

// UploadFileHandler 用户端上传文件
func UploadFileHandler(c *gin.Context) {
	userUUIDBytes, _ := c.Get("uuid")
	userUUID := userUUIDBytes.(int64)
	sha1Front := c.GetHeader("sha1")
	fileUUID, err := snowflake.GetId(1, 1)
	if err != nil {
		logger.Zap().Error(err.Error())
		resp.SendResponse(c, http.StatusInternalServerError, resp.FileUploadFailedCode, nil)
		return
	}
	fileHeader, err := c.FormFile("file")
	if err != nil {
		logger.Zap().Error(errorReadFileFromRequest, zap.Error(err))
		resp.SendResponse(c, http.StatusBadRequest, resp.FileUploadFailedCode, nil)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		logger.Zap().Error(errorOpenFile, zap.Error(err))
		resp.SendResponse(c, http.StatusBadRequest, resp.FileUploadFailedCode, nil)
		return
	}

	defer file.Close()

	filePath := uploadPath + fileHeader.Filename
	newFile, err := os.Create(filePath)
	if err != nil {
		logger.Zap().Error(errorCreateFile, zap.Error(err))
		resp.SendResponse(c, http.StatusBadRequest, resp.FileUploadFailedCode, nil)
		return
	}

	defer newFile.Close()

	io.Copy(newFile, file)

	fileBytes, err := utils.GetFileBytes(filePath)
	if err != nil {
		logger.Zap().Error(err.Error())
		resp.SendResponse(c, http.StatusInternalServerError, resp.FileUploadFailedCode, nil)
		return
	}

	// 计算文件哈希值
	sha1Backend := utils.EncryptBytesSHA1(fileBytes)

	if sha1Front != sha1Backend { // 传输过程中文件被修改
		logger.Zap().Error("file sha1 not equal", zap.String("sha1Front", sha1Front), zap.String("sha1Backend", sha1Backend))
		resp.SendResponse(c, http.StatusBadRequest, resp.FileUploadFailedCode, nil)
		return
	}

	// 数据持久化
	fileModel := &model.FileModel{
		UUID:   fileUUID,
		UserId: userUUID,
		SHA1:   sha1Backend,
		Name:   utils.CreateFileName("test"),
		Size:   fileHeader.Size,
		Path:   filePath,
	}
	userFileModel := &model.UserFileModel{
		UserId:   userUUID,
		FileId:   fileUUID,
		FileName: fileHeader.Filename,
	}

	err = file_store.Upload(fileModel, userFileModel)
	if err != nil {
		logger.Zap().Error(err.Error())
		resp.SendResponse(c, http.StatusInternalServerError, resp.FileUploadFailedCode, nil)
		return
	}

	logger.Zap().Info(successUploadFile, zap.String("file", filePath))
	resp.SendResponse(c, http.StatusOK, resp.SuccessCode, nil)
}

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

// FileChunkHandler 秒传接口
func FileChunkHandler(c *gin.Context) {
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

// DownloadFileHandler 用户端下载文件
func DownloadFileHandler(c *gin.Context) {

}

// DeleteFileHandler 用户端删除文件
func DeleteFileHandler(c *gin.Context) {

}

// RenameFileHandler 用户端重命名文件
func RenameFileHandler(c *gin.Context) {

}
