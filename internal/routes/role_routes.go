package routes

import (
	"authentication/config"
	"authentication/internal/controller"
	"github.com/gin-gonic/gin"
)

func RoleRoutes(r *gin.Engine, middleware config.Middleware, roleController controller.RoleController) {
	protected := r.Group("/v1/role")
	protected.Use(middleware.AdminMiddleware.Handler())
	{
		protected.POST("/add", roleController.AddRole)
		protected.PUT("/update/:id", roleController.UpdateRole)
		protected.GET("", roleController.GetListRole)
		protected.GET("/users", roleController.GetListRoleUsers)
		protected.GET("/users/:id", roleController.GetListUserRole)
		protected.GET("/:id", roleController.GetRoleById)
		protected.DELETE("/:id", roleController.DeleteRoleById)
	}
}
