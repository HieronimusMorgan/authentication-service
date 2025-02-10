package main

import (
	"authentication/config"
	"authentication/internal/routes"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	serverConfig, err := config.NewServerConfig()
	if err != nil {
		log.Fatalf("❌ Failed to initialize server: %v", err)
	}

	// Ensure database connection closes when the server shuts down
	defer func() {
		sqlDB, _ := serverConfig.DB.DB()
		sqlDB.Close()
		log.Println("✅ Database connection closed")
	}()

	// Start server config (Ensure everything is ready)
	if err := serverConfig.Start(); err != nil {
		log.Fatalf("❌ Error starting server: %v", err)
	}

	// Initialize Router
	r := gin.Default()

	// Register routes
	routes.ResourceRoutes(r, serverConfig.Handler.ResourceHandler)
	routes.AuthRoutes(r, serverConfig.Middleware, serverConfig.Handler.AuthHandler)
	routes.RoleRoutes(r, serverConfig.Handler.RoleHandler)

	// Run server
	log.Println("Starting server on :8080")
	err = r.Run(":8080")
	if err != nil {
		return
	}
}
