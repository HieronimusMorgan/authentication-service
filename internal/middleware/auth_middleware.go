package middleware

import (
	"authentication/internal/models"
	"authentication/internal/utils"
	"authentication/package/response"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.SendResponse(c, http.StatusUnauthorized, "Missing token", nil, "Missing token")
			c.Abort()
			return
		}

		var session models.UserSession
		if err := db.Where("session_token = ? AND is_active = ?", token, true).First(&session).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or inactive token"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			c.Abort()
			return
		}

		var resourceExists bool
		err := db.Table("authentication.user_roles").
			Select("COUNT(*) > 0 AS resource_exists").
			Joins("JOIN authentication.roles ON user_roles.role_id = roles.role_id").
			Joins("JOIN authentication.role_resources ON roles.role_id = role_resources.role_id").
			Joins("JOIN authentication.resources ON role_resources.resource_id = resources.resource_id").
			Where("user_roles.user_id = ?", session.UserID).
			Where("resources.name LIKE ?", "Auth").
			Scan(&resourceExists).Error

		if err != nil {
			// Handle error
			response.SendResponse(c, http.StatusUnauthorized, "Admin Access", nil, err.Error())
			c.Abort()
			return
			log.Println("Error executing query:", err)
		} else if resourceExists {
			log.Println("Resource exists for the user")
		} else {
			log.Println("Resource does not exist for the user")
		}
		// Check if the session has expired
		if time.Now().After(session.ExpiresAt) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		// Attach user information to the context for use in subsequent handlers
		c.Set("user_id", session.UserID)
		c.Set("session_id", session.UserSessionID)

		_, err = utils.ValidateTokenAdmin(session.SessionToken)
		if err != nil {
			response.SendResponse(c, http.StatusUnauthorized, "Admin Access", nil, err.Error())
			c.Abort()
			return
		}

		c.Next()
	}
}
