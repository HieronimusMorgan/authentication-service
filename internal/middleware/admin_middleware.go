package middleware

import (
	"authentication/internal/utils"
	"authentication/package/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AdminMiddleware defines the contract for authentication middleware
type AdminMiddleware interface {
	Handler() gin.HandlerFunc
}

// adminMiddleware is the struct that implements AdminMiddleware
type adminMiddleware struct {
	JWTService utils.JWTService
}

// NewAdminMiddleware initializes authentication middleware
func NewAdminMiddleware(jwtService utils.JWTService) AdminMiddleware {
	return adminMiddleware{
		JWTService: jwtService,
	}
}

// Handler returns a middleware function for JWT validation
func (a adminMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.SendResponse(c, http.StatusUnauthorized, "Missing token", nil, "Authorization header is required")
			c.Abort()
			return
		}

		_, err := a.JWTService.ValidateTokenAdmin(token)
		if err != nil {
			response.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil, err.Error())
			c.Abort()
			return
		}

		tokenClaims, err := a.JWTService.ExtractClaims(token)
		if err != nil {
			response.SendResponse(c, http.StatusUnauthorized, "Invalid token claims", nil, err.Error())
			c.Abort()
			return
		}

		c.Set("token", tokenClaims)

		c.Next()
	}
}
