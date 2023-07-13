package router

import (
	"onlinedisk-backend/handler"

	"github.com/gin-gonic/gin"
)

// SetupRouter 路由设置
func SetupRouter(r *gin.Engine) {
	r.Static("/static", "./static")

	fileController := r.Group("/file")
	fileController.POST("/upload", handler.UploadFileHandler)
}
