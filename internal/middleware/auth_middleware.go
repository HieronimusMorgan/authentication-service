package middleware

import (
	"authentication/internal/utils"
	"authentication/package/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.SendResponse(c, http.StatusUnauthorized, "Missing token", nil, "Missing token")
			c.Abort()
			return
		}

		_, err := utils.ValidateTokenAdmin(token)
		if err != nil {
			response.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil, err.Error())
			c.Abort()
			return
		}

		c.Next()
	}
}
