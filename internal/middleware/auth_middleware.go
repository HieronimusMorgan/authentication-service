package middleware

import (
	"Authentication/internal/utils"
	"Authentication/package/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Internal-Token")
		if token == "" {
			response.SendResponse(c, http.StatusUnauthorized, "Missing token", nil, "Missing token")
			c.Abort()
			return
		}

		_, err := utils.ValidateInternalToken(token)
		if err != nil {
			response.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil, "Invalid token")
			c.Abort()
			return
		}

		c.Next()
	}
}
