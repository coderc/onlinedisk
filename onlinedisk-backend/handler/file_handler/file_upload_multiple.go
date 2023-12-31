package file_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"onlinedisk-backend/pkg/file_store"
	"strconv"

	"github.com/coderc/onlinedisk-util/model"
	rabbitmq "github.com/coderc/onlinedisk-util/rabbitmq"
	rdb "github.com/coderc/onlinedisk-util/redis"
	"github.com/coderc/onlinedisk-util/snowflake"
	"github.com/coderc/onlinedisk-util/utils"
	"go.uber.org/zap"

	"os"

	"github.com/coderc/onlinedisk-util/logger"
	"github.com/gin-gonic/gin"

	requestUtil "github.com/coderc/onlinedisk-util/request"
	responseUtil "github.com/coderc/onlinedisk-util/response"
)

// FileUploadMultipleInitHandler 初始化分块上传
func FileUploadMultipleInitHandler(c *gin.Context) {
	var info requestUtil.MultipleInfo
	if err := c.BindJSON(&info); err != nil {
		logger.Zap().Error(err.Error())
		responseUtil.SendResponse(c, http.StatusBadRequest, responseUtil.FileUploadMultipleInitFailedCode, nil)
		return
	}

	// 将信息存入redis
	rdb.GetConn().HSet(context.Background(),
		info.UploadId, // hmap key
		"fileName", info.FileName,
		"fileSize", info.FileSize,
		"chunkSize", info.ChunkSize,
		"chunkCount", info.ChunkCount,
		"fileSHA1", info.FileSHA1,
	)

	logger.Zap().Info("upload multiple init success", zap.Any("info", info))
	responseUtil.SendResponse(c, http.StatusOK, responseUtil.SuccessCode, nil)
}

// FileUploadMultipleChunkHandler 上传单个分块
func FileUploadMultipleChunkHandler(c *gin.Context) {
	fileHandler, err := c.FormFile("file")
	if err != nil {
		logger.Zap().Error(err.Error())
		responseUtil.SendResponse(c, http.StatusBadRequest, responseUtil.FileUploadMultipleFailedCode, nil)
		return
	}

	uploadId := c.PostForm("uploadId")
	chunkSizeStr := c.PostForm("chunkSize")
	currentChunk := c.PostForm("currentChunk")
	fileName := c.PostForm("fileName")
	fileSize := c.PostForm("fileSize")
	fileSHA1 := c.PostForm("fileSHA1")

	if !utils.CheckStrsIsEmpty(uploadId, chunkSizeStr, currentChunk, fileName, fileSize, fileSHA1) {
		logger.Zap().Error("params is invalid")
		responseUtil.SendResponse(c, http.StatusBadRequest, responseUtil.FileUploadMultipleFailedCode, nil)
		return
	}

	chunkSize, err := strconv.ParseInt(chunkSizeStr, 10, 64)
	if err != nil {
		logger.Zap().Error(err.Error())
		responseUtil.SendResponse(c, http.StatusBadRequest, responseUtil.FileUploadMultipleFailedCode, nil)
		return
	}

	// 将chunk存入本地
	file, err := fileHandler.Open()
	if err != nil {
		logger.Zap().Error(err.Error())
		responseUtil.SendResponse(c, http.StatusBadRequest, responseUtil.FileUploadMultipleFailedCode, nil)
		return
	}

	defer file.Close()

	filePath := uploadPath + uploadId + currentChunk
	newFile, err := os.Create(filePath)
	if err != nil {
		logger.Zap().Error(err.Error())
		responseUtil.SendResponse(c, http.StatusBadRequest, responseUtil.FileUploadMultipleFailedCode, nil)
		return
	}

	defer newFile.Close()

	io.Copy(newFile, file)

	// 将chunk信息存入redis
	rdb.GetConn().HSet(context.Background(),
		uploadId, // hmap key
		"currentChunk"+currentChunk, currentChunk,
		"chunkSize"+currentChunk, chunkSize,
	)

	logger.Zap().Info("upload chunk success", zap.String("uploadId", uploadId), zap.String("currentChunk", currentChunk))
	responseUtil.SendResponse(c, http.StatusOK, responseUtil.SuccessCode, nil)
}

// FileUploadMultipleChunkCheckHandler 检查是否存在该chunk
func FileUploadMultipleChunkCheckHandler(c *gin.Context) {
	uploadId := c.Query("uploadId") // redis key
	currentChunk := c.Query("currentChunk")

	// 检查redis中是否存在该chunk
	chunkNumInRedis, err := rdb.GetConn().HGet(context.Background(),
		uploadId,
		"currentChunk"+currentChunk).Result()

	if err != nil {
		logger.Zap().Warn(err.Error())
		responseUtil.SendResponse(c, http.StatusForbidden, responseUtil.FileUploadMultipleFailedCode, nil)
		return
	}

	if chunkNumInRedis == currentChunk { // 存在该chunk
		logger.Zap().Info("chunk exist", zap.String("uploadId", uploadId), zap.String("currentChunk", currentChunk))
		responseUtil.SendResponse(c, http.StatusOK, responseUtil.SuccessCode, nil)
		return
	}

	logger.Zap().Error("chunk not exist")
	responseUtil.SendResponse(c, http.StatusForbidden, responseUtil.FileUploadMultipleNilChunkCode, nil)
	return
}

// FileUploadMultipleMergeHandler 合并分块 并删除redis数据, 上传db数据
func FileUploadMultipleMergeHandler(c *gin.Context) {
	var info requestUtil.MultipleInfo
	if err := c.BindJSON(&info); err != nil {
		logger.Zap().Error(err.Error())
		responseUtil.SendResponse(c, http.StatusBadRequest, responseUtil.FileUploadMultipleInitFailedCode, nil)
		return
	}
	userUUIDAny, _ := c.Get("uuid")
	userUUID := userUUIDAny.(int64)
	fileUUID, err := snowflake.GetId(1, 1)
	if err != nil {
		logger.Zap().Error(err.Error())
		responseUtil.SendResponse(c, http.StatusInternalServerError, responseUtil.FileUploadMultipleFailedCode, nil)
		return
	}

	// 合并分块
	err = mergeChunk(uploadPath, info.UploadId, int(info.ChunkCount))
	if err != nil {
		logger.Zap().Error(err.Error())
		responseUtil.SendResponse(c, http.StatusInternalServerError, responseUtil.FileUploadMultipleFailedCode, nil)
		return
	}

	// 删除redis数据
	rdb.GetConn().Del(context.TODO(), info.UploadId)

	// 数据持久化
	fileModel := &model.FileModel{
		UUID:   fileUUID,
		UserId: userUUID,
		SHA1:   info.FileSHA1,
		Name:   utils.CreateFileName("test"),
		Size:   info.FileSize,
		Path:   fmt.Sprintf("%s%s", uploadPath, info.UploadId),
	}
	userFileModel := &model.UserFileModel{
		UserId:   userUUID,
		FileId:   fileUUID,
		FileName: info.FileName,
	}

	err = file_store.Upload(fileModel, userFileModel)
	if err != nil {
		logger.Zap().Error(err.Error())
		responseUtil.SendResponse(c, http.StatusInternalServerError, responseUtil.FileUploadFailedCode, nil)
		return
	}

	// 发送消息到rabbitmq
	jsonBytes, _ := json.Marshal(fileModel)
	_ = rabbitmq.Send("file-upload-oss", jsonBytes)
	logger.Zap().Info("multiple upload success", zap.String("uploadId", info.UploadId))
	responseUtil.SendResponse(c, http.StatusOK, responseUtil.SuccessCode, nil)
}

// mergeChunk 合并分块
func mergeChunk(path, uploadId string, chunkCount int) error {
	// 创建一个新文件
	newFile, err := os.Create(path + uploadId)
	if err != nil {
		logger.Zap().Error(err.Error())
		return err
	}

	defer newFile.Close()

	// 读取分块文件
	for i := 0; i < chunkCount; i++ {
		currentChunkPath := fmt.Sprintf("%s%s%d", path, uploadId, i)
		currentChunkFile, err := os.Open(currentChunkPath)
		if err != nil {
			logger.Zap().Error(err.Error())
			return err
		}

		defer func() {
			currentChunkFile.Close()
			os.Remove(currentChunkPath)
		}()

		// 将分块文件写入新文件
		io.Copy(newFile, currentChunkFile)
	}
	return nil
}
