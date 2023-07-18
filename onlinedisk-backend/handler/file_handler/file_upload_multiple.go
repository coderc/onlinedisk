package file_handler

import (
	"context"
	"net/http"
	resp "onlinedisk-backend/response_builder"

	rdb "github.com/coderc/onlinedisk-util/redis"
	"go.uber.org/zap"

	"github.com/coderc/onlinedisk-util/logger"
	"github.com/gin-gonic/gin"
)

type MultipleInitInfo struct {
	FileName   string `json:"fileName"`
	FileSize   int64  `json:"fileSize"`
	ChunkSize  int64  `json:"chunkSize"`
	ChunkCount int64  `json:"chunkCount"`
	FileSHA1   string `json:"fileSHA1"`
	UploadId   string `json:"uploadId"`
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

}
