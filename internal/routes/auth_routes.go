package routes

import (
	"authentication/config"
	"authentication/internal/controller"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine, middleware config.Middleware, authHandler controller.AuthHController) {

	// Public routes: No middleware applied
	public := r.Group("/v1")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
	}

	protected := r.Group("/v1")
	protected.Use(middleware.AuthMiddleware.Handler())
	{
		protected.GET("/profile", authHandler.GetProfile)
		protected.POST("/change-password", authHandler.ChangePassword)
		public.GET("/logout", authHandler.Logout)
	}

	admin := r.Group("/v1")
	admin.Use(middleware.AdminMiddleware.Handler())
	{
		admin.GET("/users", authHandler.GetListUser)
		admin.POST("/user/update-role/:id", authHandler.UpdateRole)
		admin.DELETE("/user/:id", authHandler.DeleteUser)
	}
}
