package routes

import (
	"Authentication/internal/handler"
	"Authentication/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRoutes(r *gin.Engine, db *gorm.DB) {
	// Initialize Handlers
	authHandler := handler.NewAuthHandler(db)

	// Public Routes
	public := r.Group("/auth/v1")
	{
		public.POST("/register/internal-token", authHandler.RegisterInternalToken)
		public.Use(middleware.AuthMiddleware())
		{
			public.POST("/register", authHandler.Register)
			public.POST("/login", authHandler.Login)
			//public.POST("/refresh", authHandler.RefreshToken)
			public.GET("/profile", authHandler.GetProfile)

		}
	}

	// Default health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
