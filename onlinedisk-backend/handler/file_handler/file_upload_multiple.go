package file_handler

import (
	"context"
	"io"
	"net/http"
	resp "onlinedisk-backend/response_builder"
	"strconv"

	rdb "github.com/coderc/onlinedisk-util/redis"
	"github.com/coderc/onlinedisk-util/utils"
	"go.uber.org/zap"

	"os"

	"github.com/coderc/onlinedisk-util/logger"
	"github.com/gin-gonic/gin"
)

type MultipleInitInfo struct {
	UploadId   string `json:"uploadId"`
	FileName   string `json:"fileName"`
	FileSize   int64  `json:"fileSize"`
	ChunkSize  int64  `json:"chunkSize"`
	ChunkCount int64  `json:"chunkCount"`
	FileSHA1   string `json:"fileSHA1"`
}

func FileUploadMultipleInitHandler(c *gin.Context) {
	var info MultipleInitInfo
	if err := c.BindJSON(&info); err != nil {
		logger.Zap().Error(err.Error())
		resp.SendResponse(c, http.StatusBadRequest, resp.FileUploadMultipleInitFailedCode, nil)
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
	resp.SendResponse(c, http.StatusOK, resp.SuccessCode, nil)
}

func FileUploadMultipleHandler(c *gin.Context) {
	fileHandler, err := c.FormFile("file")
	if err != nil {
		logger.Zap().Error(err.Error())
		resp.SendResponse(c, http.StatusBadRequest, resp.FileUploadMultipleFailedCode, nil)
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
		resp.SendResponse(c, http.StatusBadRequest, resp.FileUploadMultipleFailedCode, nil)
		return
	}

	chunkSize, err := strconv.ParseInt(chunkSizeStr, 10, 64)
	if err != nil {
		logger.Zap().Error(err.Error())
		resp.SendResponse(c, http.StatusBadRequest, resp.FileUploadMultipleFailedCode, nil)
		return
	}

	// 将chunk存入本地
	file, err := fileHandler.Open()
	if err != nil {
		logger.Zap().Error(err.Error())
		resp.SendResponse(c, http.StatusBadRequest, resp.FileUploadMultipleFailedCode, nil)
		return
	}

	defer file.Close()

	filePath := uploadPath + uploadId + currentChunk
	newFile, err := os.Create(filePath)
	if err != nil {
		logger.Zap().Error(err.Error())
		resp.SendResponse(c, http.StatusBadRequest, resp.FileUploadMultipleFailedCode, nil)
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
	resp.SendResponse(c, http.StatusOK, resp.SuccessCode, nil)
}
