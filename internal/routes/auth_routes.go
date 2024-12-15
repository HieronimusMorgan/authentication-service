package routes

import (
	"authentication/internal/handler"
	"authentication/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRoutes(r *gin.Engine, db *gorm.DB) {
	// Initialize Handlers
	authHandler := handler.NewAuthHandler(db)

	// Public Routes
	public := r.Group("/auth/v1")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
		public.GET("/profile", authHandler.GetProfile)
	}
	public.Use(middleware.AuthMiddleware())
	{
		public.DELETE("/delete-user/:id", authHandler.DeleteUser)
	}

}
