package routes

import (
	"authentication/config"
	"authentication/internal/controller"
	"github.com/gin-gonic/gin"
)

func ResourceRoutes(r *gin.Engine, middleware config.Middleware, resourceController controller.ResourceController) {
	protected := r.Group("/v1/resources")
	protected.Use(middleware.AdminMiddleware.Handler())
	{
		protected.POST("/add", resourceController.AddResource)
		protected.POST("/update/:id", resourceController.UpdateResource)
		protected.POST("/assign-user-resources", resourceController.AssignUserResource)
		protected.POST("/remove-user-resources", resourceController.RemoveAssignUserResource)
		protected.GET("", resourceController.GetResources)
		protected.GET("/users", resourceController.GetUserResources)
		protected.GET("/:id", resourceController.GetResourcesById)
		protected.GET("/user/:id", resourceController.GetResourceUserById)
		protected.DELETE("/:id", resourceController.DeleteResourceById)
	}
}
