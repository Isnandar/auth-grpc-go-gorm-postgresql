package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
)

var RedisClient *redis.Client

func InitRedis() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	redisAddr := os.Getenv("REDIS_ADDRESS")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := os.Getenv("REDIS_DB")

	if redisAddr == "" {
		redisAddr = "localhost"
	}
	if redisPort == "" {
		redisPort = "6379"
	}
	if redisPassword == "" {
		redisPassword = ""
	}
	if redisDB == "" {
		redisDB = "0"
	}

	redisFullAddr := fmt.Sprintf("%s:%s", redisAddr, redisPort)

	db, err := strconv.Atoi(redisDB)
	if err != nil {
		log.Fatalf("Error parsing REDIS_DB: %v", err)
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisFullAddr,
		Password: redisPassword,
		DB:       db,
	})

	_, err = RedisClient.Ping().Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	fmt.Println("Connected to Redis")
}

func SetRedisValue(key string, value interface{}, expiration time.Duration) error {
	var data []byte
	var err error

	switch v := value.(type) {
	case string:
		data = []byte(v)
	default:
		data, err = json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal value to JSON: %w", err)
		}
	}

	return RedisClient.Set(key, data, expiration).Err()
}

func GetRedisValue(key string) (string, error) {
	return RedisClient.Get(key).Result()
}

func DeleteRedisValue(key string) error {
	return RedisClient.Del(key).Err()
}
