package router

import (
	fileHandler "onlinedisk-backend/handler/file_handler"
	userHandler "onlinedisk-backend/handler/user_handler"
	"onlinedisk-backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 路由设置
func SetupRouter(r *gin.Engine) {
	r.Static("/static", "./static")

	apiRouter := r.Group("/api")

	v1Router := apiRouter.Group("/v1")
	{
		fileRouter := v1Router.Group("/file")
		fileRouter.Use(middleware.Jwt())
		{
			fileRouter.POST("/upload", fileHandler.FileUploadHandler)
			fileRouter.POST("/chunk", fileHandler.FileSecondPassHandler)
			fileRouter.GET("/list", fileHandler.FileListHandler)
			fileRouter.GET("/download", fileHandler.FileDownloadHandler)
			fileRouter.POST("/delete", fileHandler.FileDeleteHandler)
		}

		userRouter := v1Router.Group("/user")
		{
			userRouter.POST("/signup", userHandler.SignupHandler)
			userRouter.POST("/signin", userHandler.SigninHandler)
		}
	}
}
