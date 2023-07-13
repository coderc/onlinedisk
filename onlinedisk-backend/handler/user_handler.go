package handler

import (
	requestInfo "github.com/coderc/onlinedisk-util/request"
	"github.com/gin-gonic/gin"
)

func SignUpHandler(c *gin.Context) {
	var userInfo requestInfo.RequestUserInfo
	if err := c.BindJSON(&userInfo); err != nil {
		c.JSON(400, gin.H{
			"msg": "参数错误",
		})
		return
	}
}
