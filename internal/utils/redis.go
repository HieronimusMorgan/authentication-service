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
	SaveData(key, clientID string, data interface{}) error
	SaveDataExpired(key, clientID string, exp float32, data interface{}) error
	GetData(key, clientID string, target interface{}) error
	DeleteData(key, clientID string) error
	GetToken(clientID string) (string, error)
	DeleteToken(clientID string) error
}

// redisService implements RedisService
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
func (r redisService) SaveData(key, clientID string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}
	return r.Client.Set(r.Ctx, key+":"+clientID, jsonData, 0).Err()
}

// SaveDataExpired stores data in Redis with expiration
func (r redisService) SaveDataExpired(key, clientID string, exp float32, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}
	return r.Client.Set(r.Ctx, key+":"+clientID, jsonData, time.Duration(exp)*time.Minute).Err()
}

// GetData retrieves and unmarshals data from Redis
func (r redisService) GetData(key, clientID string, target interface{}) error {
	jsonData, err := r.Client.Get(r.Ctx, key+":"+clientID).Result()
	if errors.Is(err, redis.Nil) {
		return fmt.Errorf("no data found for key: %s", key+":"+clientID)
	} else if err != nil {
		return fmt.Errorf("failed to get data: %v", err)
	}
	return json.Unmarshal([]byte(jsonData), target)
}

// DeleteData removes a key from Redis
func (r redisService) DeleteData(key, clientID string) error {
	return r.Client.Del(r.Ctx, key+":"+clientID).Err()
}

// generateRedisKey creates a formatted key for token storage
func generateRedisKey(clientID string) string {
	return "token:" + clientID
}

// GetToken retrieves a stored token from Redis
func (r redisService) GetToken(clientID string) (string, error) {
	token, err := r.Client.Get(r.Ctx, generateRedisKey(clientID)).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	return token, err
}

// DeleteToken removes a stored token from Redis
func (r redisService) DeleteToken(clientID string) error {
	return r.Client.Del(r.Ctx, generateRedisKey(clientID)).Err()
}

// GetUserRedis retrieves a user from Redis
func GetUserRedis(redis RedisService, key, clientID string) (*models.Users, error) {
	var user models.Users
	if err := redis.GetData(key, clientID, &user); err != nil {
		return nil, err
	}
	return &user, nil
}
