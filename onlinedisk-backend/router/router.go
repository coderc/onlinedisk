package router

import (
	"onlinedisk-backend/handler"
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
			fileRouter.POST("/upload", handler.UploadFileHandler)
		}

		userRouter := v1Router.Group("/user")
		{
			userRouter.POST("/signup", handler.SignupHandler)
			userRouter.POST("/signin", handler.SigninHandler)
		}
	}
}
