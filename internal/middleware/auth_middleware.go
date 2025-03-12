package middleware

import (
	"authentication/internal/utils"
	"authentication/package/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthMiddleware defines the contract for authentication middleware
type AuthMiddleware interface {
	Handler() gin.HandlerFunc
}

// authMiddleware is the struct that implements AuthMiddleware
type authMiddleware struct {
	JWTService utils.JWTService
}

// NewAuthMiddleware initializes authentication middleware
func NewAuthMiddleware(jwtService utils.JWTService) AuthMiddleware {
	return authMiddleware{
		JWTService: jwtService,
	}
}

func (a authMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.SendResponse(c, http.StatusUnauthorized, "Missing token", nil, "Authorization header is required")
			c.Abort()
			return
		}

		_, err := a.JWTService.ValidateToken(token)
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
