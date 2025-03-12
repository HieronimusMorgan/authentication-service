package main

import (
	"authentication/config"
	"authentication/internal/routes"
	"log"
)

func main() {
	serverConfig, err := config.NewServerConfig()
	if err != nil {
		log.Fatalf("❌ Failed to initialize server: %v", err)
	}

	defer func() {
		sqlDB, _ := serverConfig.DB.DB()
		sqlDB.Close()
		log.Println("✅ Database connection closed")
	}()

	if err := serverConfig.Start(); err != nil {
		log.Fatalf("❌ Error starting server: %v", err)
	}

	engine := serverConfig.Gin

	// Register routes
	routes.ResourceRoutes(engine, serverConfig.Middleware, serverConfig.Controller.ResourceController)
	routes.AuthRoutes(engine, serverConfig.Middleware, serverConfig.Controller.AuthController)
	routes.RoleRoutes(engine, serverConfig.Middleware, serverConfig.Controller.RoleController)
	routes.UserRoutes(engine, serverConfig.Middleware, serverConfig.Controller.UserController)

	// Run server
	log.Println("Starting server on :8080")
	err = engine.Run(":8080")
	if err != nil {
		return
	}
}
