package database

import (
	"Authentication/config"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
)

var Db *gorm.DB

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
	sqlDB.Close()
}

func schemaNamingStrategy(schemaName string) schema.NamingStrategy {
	return schema.NamingStrategy{
		TablePrefix: schemaName + ".", // Use the schema as a prefix
	}
}

func RunMigrations(dsn string, migrationsPath string) {
	m, err := migrate.New(migrationsPath, dsn)
	if err != nil {
		log.Fatalf("Failed to initialize migrations: %v", err)
	}

	// Run migrations up
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations applied successfully")
}
