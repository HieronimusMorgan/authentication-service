package utils

import (
	"authentication/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

// RedisService defines the contract for Redis operations
type RedisService interface {
	SaveData(key string, clientID string, data interface{}) error
	SaveDataExpired(key string, clientID string, exp int, data interface{}) error
	GetData(key string, clientID string, target interface{}) error
	DeleteData(key string, clientID string) error
	GetToken(clientID string) (string, error)
	DeleteToken(clientID string) error
}

// RedisServiceImpl is the struct that implements RedisService
type redisService struct {
	Client redis.Client
	Ctx    context.Context
}

// NewRedisService initializes Redis client
func NewRedisService(client redis.Client) RedisService {
	return redisService{
		Client: client,
		Ctx:    context.Background(),
	}
}

// SaveData stores data in Redis
func (r redisService) SaveData(key string, clientID string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %v", err)
	}
	redisKey := key + ":" + clientID
	err = r.Client.Set(r.Ctx, redisKey, jsonData, 0).Err()
	return err
}

// SaveDataExpired stores data in Redis with an expiration time
func (r redisService) SaveDataExpired(key string, requestID string, exp int, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %v", err)
	}
	redisKey := key + ":" + requestID
	err = r.Client.Set(r.Ctx, redisKey, jsonData, time.Duration(exp)*time.Minute).Err()
	return err
}

// GetData retrieves data from Redis and unmarshals it into target
func (r redisService) GetData(key string, clientID string, target interface{}) error {
	redisKey := key + ":" + clientID
	jsonData, err := r.Client.Get(r.Ctx, redisKey).Result()
	if errors.Is(err, redis.Nil) {
		return fmt.Errorf("no data found for key: %s", redisKey)
	} else if err != nil {
		return fmt.Errorf("failed to get data from Redis: %v", err)
	}

	err = json.Unmarshal([]byte(jsonData), target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON data: %v", err)
	}
	return nil
}

// DeleteData removes a key from Redis
func (r redisService) DeleteData(key string, clientID string) error {
	redisKey := key + ":" + clientID
	err := r.Client.Del(r.Ctx, redisKey).Err()
	return err
}

// GenerateRedisKey creates a formatted key for token storage
func generateRedisKey(clientID string) string {
	return "token:" + clientID
}

// GetToken retrieves a stored token from Redis
func (r redisService) GetToken(clientID string) (string, error) {
	redisKey := generateRedisKey(clientID)
	token, err := r.Client.Get(r.Ctx, redisKey).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil // Token not found
	} else if err != nil {
		return "", err // Other errors
	}
	return token, nil
}

// DeleteToken removes a stored token from Redis
func (r redisService) DeleteToken(clientID string) error {
	redisKey := generateRedisKey(clientID)
	err := r.Client.Del(r.Ctx, redisKey).Err()
	return err
}

func GetUserRedis(redis RedisService, key string, clientID string) (*models.Users, error) {
	var u = &models.Users{}
	err := redis.GetData(key, clientID, u)
	if err != nil {
		return nil, err
	}
	return u, nil
}
