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
			c.JSON(http.StatusUnauthorized, nil)
			return
		}
		uuid, deadline, err := jwt.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, nil)
			return
		}
		if deadline < time.Now().Unix() {
			c.JSON(http.StatusUnauthorized, nil)
			return
		}
		c.Set("uuid", uuid)
		c.Next()
	}
}
