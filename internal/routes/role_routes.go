package routes

import (
	"authentication/internal/controller"
	"github.com/gin-gonic/gin"
)

func RoleRoutes(r *gin.Engine, roleHandler controller.RoleController) {

	protected := r.Group("/v1/role")
	//protected.Use(middleware.AuthMiddleware(roleService))
	{
		protected.POST("/add", roleHandler.AddRole)
		protected.PUT("/update/:id", roleHandler.UpdateRole)
		protected.GET("", roleHandler.GetListRole)
		protected.GET("/users", roleHandler.GetListRoleUsers)
		protected.GET("/:id", roleHandler.GetRoleById)
		protected.DELETE("/:id", roleHandler.DeleteRoleById)
	}
}
