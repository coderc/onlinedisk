package main

import (
	"fmt"
	"onlinedisk-backend/router"

	"github.com/coderc/onlinedisk-util/config"
	"github.com/coderc/onlinedisk-util/db"
	"github.com/coderc/onlinedisk-util/logger"
	"github.com/gin-gonic/gin"
)

func Init() {
	config.GetConfig().Init("./config/config.yaml")
	logger.GetLogger().Init(&config.GetConfig().LoggerConfig)
	db.GetSingleton().Init(&config.GetConfig().DBConfig)
}

func main() {
	Init()
	r := gin.New()
	router.SetupRouter(r)

	err := r.Run(fmt.Sprintf(":%d", config.GetConfig().ServerConfig.Port))
	if err != nil {
		panic(err)
	}
}
