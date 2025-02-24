package routes

import (
	"authentication/internal/controller"
	"github.com/gin-gonic/gin"
)

func ResourceRoutes(r *gin.Engine, resourceHandler controller.ResourceController) {
	// Initialize Handlers
	// Public Routes
	protected := r.Group("/v1/resources")
	//protected.Use(middleware.AuthMiddleware(resourceService))
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
