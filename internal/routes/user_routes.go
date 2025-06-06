package routes

import (
	"authentication/config"
	"authentication/internal/controller"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine, middleware config.Middleware, userController controller.UserController) {
	protected := r.Group("/v1")
	protected.Use(middleware.AuthMiddleware.Handler())
	{
		protected.GET("/profile", userController.GetProfile)
		protected.PUT("/update/profile-name/:id", userController.UpdateNameUserProfile)
		protected.POST("/update/profile-photo/:id", userController.UpdatePhotoUserProfile)
		protected.POST("/update/user-setting", userController.UpdateUserSetting)
		protected.DELETE("/delete-user/:id", userController.DeleteUser)
	}
}
