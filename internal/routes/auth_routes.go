package routes

import (
	"authentication/internal/handler"
	"authentication/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRoutes(r *gin.Engine, db *gorm.DB) {
	authHandler := handler.NewAuthHandler(db)

	// Public routes: No middleware applied
	public := r.Group("/auth/v1")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
	}

	// Protected routes: Require general middleware
	protected := r.Group("/auth/v1")
	protected.Use(middleware.Middleware())
	{
		protected.GET("/profile", authHandler.GetProfile)
		protected.POST("/change-password", authHandler.ChangePassword)
	}

	// Admin routes: Require authentication middleware
	admin := r.Group("/auth/v1")
	admin.Use(middleware.AuthMiddleware())
	{
		admin.GET("/users", authHandler.GetListUser)
		admin.POST("/user/update-role/:id", authHandler.UpdateRole)
		admin.DELETE("/user/:id", authHandler.DeleteUser)
	}
}
