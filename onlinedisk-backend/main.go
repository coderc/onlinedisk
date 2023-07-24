package main

import (
	"fmt"
	"onlinedisk-backend/router"

	userServiceClient "onlinedisk-backend/grpc/client"

	amqpConn "github.com/coderc/onlinedisk-util/rabbitmq"

	"github.com/coderc/onlinedisk-util/config"
	"github.com/coderc/onlinedisk-util/db"
	"github.com/coderc/onlinedisk-util/logger"
	"github.com/coderc/onlinedisk-util/redis"
	"github.com/gin-gonic/gin"

	ossClient "onlinedisk-backend/oss_client"
)

func Init() {
	config.GetConfig().Init("./config/config.yaml")
	logger.GetLogger().Init(&config.GetConfig().LoggerConfig)
	db.GetSingleton().Init(&config.GetConfig().DBConfig)
	amqpConn.Init(&config.GetConfig().RabbitMQConfig)
	redis.Init(&config.GetConfig().RedisConfig)
	ossClient.GetClient().Init(&config.GetConfig().OssConfig)

	userServiceClient.Init(&config.GetConfig().UserServiceConfig)
}

func main() {
	Init()
	r := gin.Default()
	router.SetupRouter(r)
	go ossClient.GetClient().UploadFile("file-upload-oss")
	err := r.Run(fmt.Sprintf(":%d", config.GetConfig().ServerConfig.Port))
	if err != nil {
		panic(err)
	}
}
