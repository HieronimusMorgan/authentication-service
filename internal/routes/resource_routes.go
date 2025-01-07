package routes

import (
	"authentication/internal/handler"
	"authentication/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ResourceRoutes(r *gin.Engine, db *gorm.DB) {
	// Initialize Handlers
	resourceHandler := handler.NewResourceHandler(db)

	// Public Routes
	protected := r.Group("/auth/v1/resources")
	protected.Use(middleware.AuthMiddleware(db))
	{
		protected.POST("/add", resourceHandler.AddResource)
		protected.POST("/update/:id", resourceHandler.UpdateResource)
		protected.POST("/assign-role", resourceHandler.AssignResourceToRole)
		protected.GET("", resourceHandler.GetResources)
		protected.GET("/roles", resourceHandler.GetResourceRoles)
		protected.GET("/:id", resourceHandler.GetResourcesById)
		protected.GET("/user/:id", resourceHandler.GetResourceUserById)
		protected.DELETE("/:id", resourceHandler.DeleteResourceById)
	}

}
