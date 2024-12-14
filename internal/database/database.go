package database

import (
	"authentication/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
)

func InitDB() *gorm.DB {
	cfg := config.LoadConfig()
	dsn := "host=" + cfg.DBHost + " user=" + cfg.DBUser + " password=" + cfg.DBPassword +
		" dbname=" + cfg.DBName + " port=" + cfg.DBPort + " sslmode=disable"

	var err error
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Info),
		NamingStrategy: schemaNamingStrategy(cfg.DBSchema), // Set the schema
	})

	//err = db.AutoMigrate(&models.User{}, &models.Role{}, &models.Resource{}, &models.RoleResource{}, &models.UserRole{}, &models.InternalToken{})
	//if err != nil {
	//	log.Fatalf("Failed to migrate: %v", err)
	//}
	//Run migrations
	//migrationsPath := "file://./migrations"
	//RunMigrations(dsn, migrationsPath)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established")
	return db
}

func CloseDB(db *gorm.DB) {
	sqlDB, _ := db.DB()
	err := sqlDB.Close()
	if err != nil {
		return
	}
}

func schemaNamingStrategy(schemaName string) schema.NamingStrategy {
	return schema.NamingStrategy{
		TablePrefix: schemaName + ".", // Use the schema as a prefix
	}
}
