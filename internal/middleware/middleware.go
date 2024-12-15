package middleware

import (
	"authentication/internal/utils"
	"authentication/package/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.SendResponse(c, http.StatusUnauthorized, "Missing token", nil, "Missing token")
			c.Abort()
			return
		}

		_, err := utils.ValidateToken(token)
		if err != nil {
			response.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil, "Invalid token")
			c.Abort()
			return
		}

		c.Next()
	}
}
