package config

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

// Config holds application-wide configurations
type Config struct {
	AppPort    string `envconfig:"APP_PORT" default:"8080"`
	JWTSecret  string `envconfig:"JWT_SECRET" default:"a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6"`
	RedisHost  string `envconfig:"REDIS_HOST" default:"localhost"`
	RedisPort  string `envconfig:"REDIS_PORT" default:"6379"`
	RedisDB    int    `envconfig:"REDIS_DB" default:"0"`
	RedisPass  string `envconfig:"REDIS_PASSWORD" default:""`
	DBHost     string `envconfig:"DB_HOST" default:"localhost"`
	DBPort     string `envconfig:"DB_PORT" default:"5432"`
	DBUser     string `envconfig:"DB_USER" default:"postgres"`
	DBPassword string `envconfig:"DB_PASSWORD" default:"admin"`
	DBName     string `envconfig:"DB_NAME" default:"postgres"`
	DBSchema   string `envconfig:"DB_SCHEMA" default:"authentication"`
	DBSSLMode  string `envconfig:"DB_SSLMODE" default:"disable"`
}

// LoadConfig loads environment variables into the Config struct
func LoadConfig() *Config {
	var cfg Config
	err := envconfig.Process("AUTH", &cfg)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	return &cfg
}

// InitDatabase initializes and returns a PostgreSQL database connection
func InitDatabase(cfg *Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Info),
		NamingStrategy: schemaNamingStrategy(cfg.DBSchema)})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("✅ Connected to PostgreSQL")
	return db
}

func schemaNamingStrategy(schemaName string) schema.NamingStrategy {
	return schema.NamingStrategy{
		TablePrefix: schemaName + ".", // Use the schema as a prefix
	}
}

// InitRedis initializes and returns a Redis client
func InitRedis(cfg *Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPass,
		DB:       cfg.RedisDB,
	})

	// Ping Redis to check the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}
	log.Println("✅ Connected to Redis")
	return rdb
}

func LoadRedisConfig() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv("REDIS_URL")
}
