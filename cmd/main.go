package main

import (
	"authentication/internal/database"
	"authentication/internal/routes"
	"authentication/internal/utils"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// Initialize Redis
	utils.InitializeRedis()

	// Initialize database
	db := database.InitDB()
	defer database.CloseDB(db)

	// Setup Gin router
	r := gin.Default()

	// Register routes
	routes.ResourceRoutes(r, db)
	routes.AuthRoutes(r, db)
	routes.RoleRoutes(r, db)

	// Run server
	log.Println("Starting server on :8080")
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
