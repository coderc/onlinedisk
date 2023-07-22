package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/coderc/onlinedisk-util/jwt"
)

func Jwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token == "" {
			token = c.Query("token")
		}
		if token == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		uuid, deadline, err := jwt.ParseToken(token)
		if err != nil || deadline < time.Now().Unix() {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Set("uuid", uuid)
		c.Next()
	}
}
