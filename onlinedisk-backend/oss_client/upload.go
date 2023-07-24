package oss_client

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/coderc/onlinedisk-util/logger"
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
				fmt.Println("[error]: " + err.Error())
			} else {
				o.uploadFile(fileModel)
			}
			fmt.Println("oss client upload file: [" + string(msg) + "]")
		}
	}
}

func (o *OssClient) uploadFile(fileModel *model.FileModel) {
	err := o.bucket.PutObjectFromFile(fileModel.SHA1, fileModel.Path)
	if err != nil {
		logger.Zap().Error(err.Error(), zap.String("fileSHA1", fileModel.SHA1))
		fmt.Println("[error]: " + err.Error())
	}

	// 删除本地文件
	_ = os.Remove(fileModel.Path)

	// 更新数据库
	// TODO
}
