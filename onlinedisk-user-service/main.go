package main

import (
	"github.com/coderc/onlinedisk-util/config"
	"github.com/coderc/onlinedisk-util/db"
	"github.com/coderc/onlinedisk-util/logger"
)

func Init() {
	config.GetConfig().Init("./config/config.yaml")
	logger.GetLogger().Init(&config.GetConfig().LoggerConfig)
	db.GetSingleton().Init(&config.GetConfig().DBConfig)
}

func main() {
	Init()
}
