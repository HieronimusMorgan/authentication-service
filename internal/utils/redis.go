package utils

import (
	"authentication/config"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitializeRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.LoadRedisConfig(),
		Password: "",
		DB:       0,
	})
	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis successfully")
}

func SaveTokenToRedis(clientID string, token string, expiration time.Duration) error {
	redisKey := generateRedisKey(clientID) // Generate a unique key for the user
	fmt.Println(redisKey)
	err := RedisClient.Set(Ctx, redisKey, token, expiration).Err()
	return err
}

func SaveDataToRedis(key string, clientID string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %v", err)
	}
	redisKey := key + ":" + clientID
	err = RedisClient.Set(Ctx, redisKey, jsonData, 0).Err()
	return err
}

func GetDataFromRedis(key string, clientID string, target **interface{}) error {
	redisKey := key + ":" + clientID
	jsonData, err := RedisClient.Get(Ctx, redisKey).Result()
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

func DeleteDataFromRedis(key string, clientID string) error {
	redisKey := key + ":" + clientID
	err := RedisClient.Del(Ctx, redisKey).Err()
	return err
}

func generateRedisKey(clientID string) string {
	return "token:" + string(clientID)
}

func GetTokenFromRedis(clientID string) (string, error) {
	redisKey := generateRedisKey(clientID)
	token, err := RedisClient.Get(Ctx, redisKey).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil // Token not found
	} else if err != nil {
		return "", err // Other errors
	}
	return token, nil
}

func DeleteTokenFromRedis(clientID string) error {
	redisKey := generateRedisKey(clientID)
	err := RedisClient.Del(Ctx, redisKey).Err()
	return err
}
