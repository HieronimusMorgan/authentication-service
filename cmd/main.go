package main

import (
	"Authentication/internal/database"
	"Authentication/internal/routes"
	"Authentication/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/schema"
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
	routes.AuthRoutes(r, db)

	// Run server
	log.Println("Starting server on :8080")
	r.Run(":8080")
}

func schemaNamingStrategy(schemaName string) schema.NamingStrategy {
	return schema.NamingStrategy{
		TablePrefix: schemaName + ".", // Use the schema as a prefix
	}
}
