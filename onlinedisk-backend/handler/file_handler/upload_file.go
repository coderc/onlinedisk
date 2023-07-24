package file_handler

import (
	"io"
	"net/http"
	"onlinedisk-backend/pkg/file_store"
	"os"

	response "github.com/coderc/onlinedisk-util/response"

	"encoding/json"

	"github.com/coderc/onlinedisk-util/logger"
	"github.com/coderc/onlinedisk-util/model"
	"github.com/coderc/onlinedisk-util/snowflake"
	"github.com/coderc/onlinedisk-util/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	rabbitmq "github.com/coderc/onlinedisk-util/rabbitmq"
)

// FileUploadHandler 用户端上传文件
func FileUploadHandler(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		logger.Zap().Error(errorReadFileFromRequest, zap.Error(err))
		response.SendResponse(c, http.StatusBadRequest, response.FileUploadFailedCode, nil)
		return
	}

	userUUIDAny, _ := c.Get("uuid")
	userUUID := userUUIDAny.(int64)
	sha1Front := c.GetHeader("sha1")
	fileUUID, err := snowflake.GetId(1, 1)
	if err != nil {
		logger.Zap().Error(err.Error())
		response.SendResponse(c, http.StatusInternalServerError, response.FileUploadFailedCode, nil)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		logger.Zap().Error(errorOpenFile, zap.Error(err))
		response.SendResponse(c, http.StatusBadRequest, response.FileUploadFailedCode, nil)
		return
	}

	defer file.Close()

	filePath := uploadPath + fileHeader.Filename
	newFile, err := os.Create(filePath)
	if err != nil {
		logger.Zap().Error(errorCreateFile, zap.Error(err))
		response.SendResponse(c, http.StatusBadRequest, response.FileUploadFailedCode, nil)
		return
	}

	defer newFile.Close()

	io.Copy(newFile, file)

	fileBytes, err := utils.GetFileBytes(filePath)
	if err != nil {
		logger.Zap().Error(err.Error())
		response.SendResponse(c, http.StatusInternalServerError, response.FileUploadFailedCode, nil)
		return
	}

	// 计算文件哈希值
	fileSHA1 := utils.EncryptBytesSHA1(fileBytes)

	if sha1Front != fileSHA1 { // 传输过程中文件被修改
		logger.Zap().Error("file sha1 not equal", zap.String("sha1Front", sha1Front), zap.String("sha1Backend", fileSHA1))
		response.SendResponse(c, http.StatusBadRequest, response.FileUploadFailedCode, nil)
		return
	}

	// 数据持久化
	fileModel := &model.FileModel{
		UUID:   fileUUID,
		UserId: userUUID,
		SHA1:   fileSHA1,
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
		response.SendResponse(c, http.StatusInternalServerError, response.FileUploadFailedCode, nil)
		return
	}

	// 将消息发送到rabbitmq
	jsonBytes, _ := json.Marshal(fileModel)
	err = rabbitmq.Send("file-upload-oss", jsonBytes)
	if err != nil {
		logger.Zap().Error(err.Error())
	}

	logger.Zap().Info(successUploadFile, zap.String("file", filePath))
	response.SendResponse(c, http.StatusOK, response.SuccessCode, nil)
}
