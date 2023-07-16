package handler

import (
	"io"
	"net/http"
	"net/url"
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
	sha1 := c.GetHeader("sha1")
	if sha1 != "" { // 尝试秒传
		fileName := c.GetHeader("filename")
		fileName, _ = url.QueryUnescape(fileName)
		fileModel := hasFile(sha1)
		if fileModel != nil && fileModel.SHA1 == sha1 { // 秒传成功
			userFileModel := &model.UserFileModel{
				UserId:   userUUID,
				FileId:   fileModel.UUID,
				FileName: fileName,
			}
			err := mapper.InsertUserFile(userFileModel)
			if err != nil {
				logger.GetLogger().Error(err.Error())
				resp.SendResponse(c, http.StatusInternalServerError, resp.FileUploadFailedCode, nil)
				return
			}
			logger.GetLogger().Info(successUploadFileChunk, zap.String("fileSHA1", sha1))
			resp.SendResponse(c, http.StatusOK, resp.SuccessCode, nil)
			return
		}
	}
	fileUUID, err := snowflake.GetId(1, 1)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		resp.SendResponse(c, http.StatusInternalServerError, resp.FileUploadFailedCode, nil)
		return
	}
	fileHeader, err := c.FormFile("file")
	if err != nil {
		logger.GetLogger().Error(errorReadFileFromRequest, zap.Error(err))
		resp.SendResponse(c, http.StatusBadRequest, resp.FileUploadFailedCode, nil)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		logger.GetLogger().Error(errorOpenFile, zap.Error(err))
		resp.SendResponse(c, http.StatusBadRequest, resp.FileUploadFailedCode, nil)
		return
	}

	defer file.Close()

	filePath := uploadPath + fileHeader.Filename
	newFile, err := os.Create(filePath)
	if err != nil {
		logger.GetLogger().Error(errorCreateFile, zap.Error(err))
		resp.SendResponse(c, http.StatusBadRequest, resp.FileUploadFailedCode, nil)
		return
	}

	defer newFile.Close()

	io.Copy(newFile, file)

	fileBytes, err := utils.GetFileBytes(filePath)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		resp.SendResponse(c, http.StatusInternalServerError, resp.FileUploadFailedCode, nil)
		return
	}

	// 计算文件哈希值
	sha1Str := utils.EncryptBytesSHA1(fileBytes)

	// 数据持久化
	fileModel := &model.FileModel{
		UUID:   fileUUID,
		UserId: userUUID,
		SHA1:   sha1Str,
		Name:   utils.CreateFileName("test"),
		Size:   fileHeader.Size,
		Path:   filePath,
	}
	userFileModel := &model.UserFileModel{
		UserId:   userUUID,
		FileId:   fileUUID,
		FileName: fileHeader.Filename,
	}
	err = mapper.InsertFile(fileModel) // table_file
	if err != nil {
		logger.GetLogger().Error(err.Error())
		resp.SendResponse(c, http.StatusInternalServerError, resp.FileUploadFailedCode, nil)
		return
	}
	err = mapper.InsertUserFile(userFileModel) // table_user_file

	if err != nil {
		logger.GetLogger().Error(err.Error())
		resp.SendResponse(c, http.StatusInternalServerError, resp.FileUploadFailedCode, nil)
		return
	}

	logger.GetLogger().Info(successUploadFile, zap.String("file", filePath))
	resp.SendResponse(c, http.StatusOK, resp.SuccessCode, nil)
}

func FileListHandler(c *gin.Context) {
	userUUIDBytes, _ := c.Get("uuid")
	userUUID := userUUIDBytes.(int64)

	userFileModels, err := mapper.QueryUserFileByUserId(userUUID)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		resp.SendResponse(c, http.StatusInternalServerError, resp.FailedSelectFileListCode, nil)
		return
	}

	logger.GetLogger().Info("select file list success",
		zap.Int64("userUUID", userUUID), zap.Int("fileCount", len(userFileModels)))
	resp.SendResponse(c, http.StatusOK, resp.SuccessCode, userFileModels)
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

func hasFile(sha1 string) *model.FileModel {
	fileModel, err := mapper.QueryFileBySHA1(sha1)
	if err != nil {
		return nil
	}
	return fileModel
}
