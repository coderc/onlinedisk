package response_builder

import (
	"github.com/gin-gonic/gin"
)

const (
	SuccessCode              = 0
	FileUploadFailedCode     = 1001
	FailedSelectFileListCode = 1002
	SignupFailedCode         = 2001
	SigninFailedCode         = 2002
)

func SendResponse(c *gin.Context, httpCode, serviceCode int, data interface{}) {
	c.JSON(httpCode, gin.H{
		"code": serviceCode,
		"data": data,
	})
}
