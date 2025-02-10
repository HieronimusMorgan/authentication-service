package routes

import (
	"authentication/internal/handler"
	"github.com/gin-gonic/gin"
)

func RoleRoutes(r *gin.Engine, roleHandler handler.RoleHandler) {

	protected := r.Group("/auth/v1/role")
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
