package oss_client

import (
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/coderc/onlinedisk-util/config"
)

var (
	once   sync.Once
	client *OssClient
)

type OssClient struct {
	client *oss.Client
	bucket *oss.Bucket
}

func GetClient() *OssClient {
	once.Do(func() {
		client = &OssClient{}
	})
	return client
}

func (o *OssClient) Init(config *config.OssConfigStruct) {
	var err error
	o.client, err = oss.New(config.Endpoint, config.AccessKeyID, config.AccessKeySecret)
	if err != nil {
		panic(err)
	}

	o.bucket, err = o.client.Bucket(config.BucketName)
	if err != nil {
		panic(err)
	}
}
