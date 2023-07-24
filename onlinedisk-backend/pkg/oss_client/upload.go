package oss_client

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/coderc/onlinedisk-util/logger"
	"github.com/coderc/onlinedisk-util/mapper"
	"github.com/coderc/onlinedisk-util/model"
	rabbitmq "github.com/coderc/onlinedisk-util/rabbitmq"
	"go.uber.org/zap"
)

func (o *OssClient) UploadFile(qName string) {
	chanMsg := make(chan []byte, 100)

	go rabbitmq.GetMsg(qName, chanMsg)

	for {
		select {
		case msg := <-chanMsg:
			fileModel := &model.FileModel{}
			err := json.Unmarshal(msg, fileModel)
			if err != nil {
				logger.Zap().Error("oss client upload file: [" + string(msg) + "]")
			} else {
				o.uploadFile(fileModel)
			}
		}
	}
}

func (o *OssClient) uploadFile(fileModel *model.FileModel) {
	defer os.Remove(fileModel.Path)
	err := o.bucket.PutObjectFromFile(strconv.FormatInt(fileModel.UUID, 10), fileModel.Path)
	if err != nil {
		logger.Zap().Error(err.Error(), zap.Int64("fileUUID", fileModel.UUID))
	} else {
		logger.Zap().Debug("oss client upload file success", zap.Int64("fileUUID", fileModel.UUID))
	}

	fileModel.Path = strconv.FormatInt(fileModel.UUID, 10)

	// 更新数据库
	err = mapper.UpdateFilePath(fileModel)
	if err != nil {
		logger.Zap().Error(err.Error(), zap.Int64("fileUUID", fileModel.UUID))
	} else {
		logger.Zap().Debug("oss client update file path success", zap.Int64("fileUUID", fileModel.UUID))
	}
}
