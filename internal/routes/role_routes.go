package routes

import (
	"authentication/internal/handler"
	"authentication/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RoleRoutes(r *gin.Engine, db *gorm.DB) {
	// Initialize Handlers
	roleHandler := handler.NewRoleHandler(db)

	protected := r.Group("/auth/v1/role")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/add", roleHandler.AddRole)
		protected.PUT("/update/:id", roleHandler.UpdateRole)
		protected.GET("", roleHandler.GetListRole)
		protected.GET("/:id", roleHandler.GetRoleById)
		protected.DELETE("/:id", roleHandler.DeleteRoleById)
	}
}
