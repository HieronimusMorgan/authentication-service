package config

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// Config holds application-wide configurations
type Config struct {
	AppPort    string `envconfig:"APP_PORT" default:"8080"`
	JWTSecret  string `envconfig:"JWT_SECRET" default:"a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6"`
	AesEncrypt string `envconfig:"AES_ENCRYPT" default:"E8gWmPq4R7YTXA9LJ0b3VNdT1KoU6C8h"`
	AesFixedIV string `envconfig:"AES_FIXED_IV" default:"D4f9s8X7a2K6mL1Q"`
	RedisHost  string `envconfig:"REDIS_HOST" default:"localhost"`
	RedisPort  string `envconfig:"REDIS_PORT" default:"6379"`
	RedisDB    int    `envconfig:"REDIS_DB" default:"0"`
	RedisPass  string `envconfig:"REDIS_PASSWORD" default:""`
	DBHost     string `envconfig:"DB_HOST" default:"localhost"`
	DBPort     string `envconfig:"DB_PORT" default:"5432"`
	DBUser     string `envconfig:"DB_USER" default:"postgres"`
	DBPassword string `envconfig:"DB_PASSWORD" default:"admin"`
	DBName     string `envconfig:"DB_NAME" default:"authentication"`
	DBSchema   string `envconfig:"DB_SCHEMA" default:"public"`
	DBSSLMode  string `envconfig:"DB_SSLMODE" default:"disable"`
	CdnUrl     string `envconfig:"CDN_URL"  default:"http://localhost:8181"`
	NatsUrl    string `envconfig:"NATS_URL" default:"nats://localhost:4222"`
}

// LoadConfig loads environment variables into the Config struct
func LoadConfig() *Config {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	logrus.WithFields(logrus.Fields{
		"AppPort":   cfg.AppPort,
		"DBHost":    cfg.DBHost,
		"DBName":    cfg.DBName,
		"RedisHost": cfg.RedisHost,
	}).Info("✅ Configuration loaded successfully")

	return &cfg
}

// InitGin initializes the Gin engine with appropriate configurations
func InitGin() *gin.Engine {
	// Set Gin mode based on environment
	if ginMode := gin.Mode(); ginMode != gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
		logrus.Warn("⚠ Running in DEBUG mode. Use `GIN_MODE=release` in production.")
	} else {
		logrus.Info("✅ Running in RELEASE mode.")
	}

	// Create a new Gin router
	engine := gin.New()

	// Middleware
	engine.Use(gin.Recovery()) // Handles panics and prevents crashes
	engine.Use(gin.Logger())   // Logs HTTP requests

	// Security Headers (Prevents Clickjacking & XSS Attacks)
	engine.Use(func(c *gin.Context) {
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Next()
	})

	logrus.Info("🚀 Gin HTTP server initialized successfully")
	return engine
}

// InitDatabase initializes and returns a PostgreSQL database connection with retry logic
func InitDatabase(cfg *Config) *gorm.DB {
	var once sync.Once
	var db *gorm.DB
	var err error
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
		)
		maxRetries := 5
		for i := 1; i <= maxRetries; i++ {
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
				Logger:         logger.Default.LogMode(logger.Info),
				NamingStrategy: schemaNamingStrategy(cfg.DBSchema),
			})
			if err == nil {
				break
			}

			logrus.WithFields(logrus.Fields{
				"attempt": i,
				"error":   err.Error(),
			}).Warn("⏳ Retrying database connection...")

			time.Sleep(2 * time.Second)
		}

		if err != nil {
			logrus.WithError(err).Fatal("❌ Failed to connect to PostgreSQL after retries")
		}

		logrus.Info("✅ Connected to PostgreSQL")
	})
	return db
}

// CloseDatabase closes the database connection properly
func CloseDatabase(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		logrus.WithError(err).Error("Failed to retrieve database instance")
		return
	}

	if err := sqlDB.Close(); err != nil {
		logrus.WithError(err).Error("Error closing database connection")
	} else {
		logrus.Info("✅ Database connection closed")
	}
}

// InitRedis initializes and returns a Redis client with retry logic
func InitRedis(cfg *Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPass,
		DB:       cfg.RedisDB,
	})

	// Connection Retry Logic
	maxRetries := 5
	for i := 1; i <= maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := rdb.Ping(ctx).Result()
		if err == nil {
			logrus.Info("✅ Connected to Redis")
			return rdb
		}

		logrus.WithFields(logrus.Fields{
			"attempt": i,
			"error":   err.Error(),
		}).Warn("⏳ Retrying Redis connection...")

		time.Sleep(2 * time.Second)
	}

	logrus.Fatal("❌ Failed to connect to Redis after retries")
	return nil
}

// CloseRedis closes the Redis connection properly
func CloseRedis(rdb *redis.Client) {
	if err := rdb.Close(); err != nil {
		logrus.WithError(err).Error("Error closing Redis connection")
	} else {
		logrus.Info("✅ Redis connection closed")
	}
}

// schemaNamingStrategy sets the schema for GORM
func schemaNamingStrategy(schemaName string) schema.NamingStrategy {
	return schema.NamingStrategy{
		TablePrefix: schemaName + ".", // Use schema as prefix
	}
}
