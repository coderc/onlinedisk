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
			fileRouter.POST("/second_pass", fileHandler.FileSecondPassHandler)
			fileRouter.GET("/list", fileHandler.FileListHandler)
			fileRouter.GET("/download", fileHandler.FileDownloadHandler)
			fileRouter.POST("/delete", fileHandler.FileDeleteHandler)

			// 分段上传: init
			fileRouter.POST("/multipart_upload/init", fileHandler.FileUploadMultipleInitHandler)
			// 分段上传: single chunk upload
			fileRouter.POST("/multipart_upload/chunk", fileHandler.FileUploadMultipleChunkHandler)
			// 断点续传上传: 检查chunk是否存在
			fileRouter.GET("/multipart_upload/check", fileHandler.FileUploadMultipleChunkCheckHandler)
			// 分段上传: merge
			fileRouter.POST("/multipart_upload/merge", fileHandler.FileUploadMultipleMergeHandler)
		}

		userRouter := v1Router.Group("/user")
		{
			userRouter.POST("/signup", userHandler.SignupHandler)
			userRouter.POST("/signin", userHandler.SigninHandler)
		}
	}
}
