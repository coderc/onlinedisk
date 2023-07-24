package oss_client

import (
	"fmt"
	"strconv"

	"github.com/coderc/onlinedisk-util/logger"
	"go.uber.org/zap"
)

func (o *OssClient) Download(fileUUID int64) (string, error) {
	// 从oss下载文件
	localPath := fmt.Sprintf("./static/file/%d", fileUUID)
	err := o.bucket.GetObjectToFile(strconv.FormatInt(fileUUID, 10), localPath)
	if err != nil {
		logger.Zap().Error(err.Error(), zap.Int64("fileUUID", fileUUID))
		return "", err
	}
	return localPath, nil
}
