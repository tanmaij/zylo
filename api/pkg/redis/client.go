package redis

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/tanmaij/zylo/config"

	"github.com/go-redis/redis/v8"
)

// Client represents a Redis client wrapper.
type Client struct {
	redisClient *redis.Client // The underlying Redis client instance.
}

// NewRedisClient new redis client instance
func NewRedisClient() (Client, error) {
	cfg := config.Instance.Redis

	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisAddr,
		MinIdleConns: cfg.MinIdleConns,
		PoolSize:     cfg.PoolSize,
		PoolTimeout:  time.Duration(cfg.PoolTimeout) * time.Second,
	})

	ctx := context.Background()
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Println("Could not connect to Redis:", err)
		return Client{}, err
	}

	log.Println("Connected to Redis:", pong)
	return Client{redisClient: rdb}, nil
}

// Get retrieves data from Redis using the specified key and unmarshals it into the provided data structure.
func (c Client) Get(ctx context.Context, key string, data interface{}) error {
	// Retrieve data as bytes from Redis using the specified key.
	dataBytes, err := c.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	// Unmarshal the data bytes into the provided data structure.
	return json.Unmarshal(dataBytes, data)
}

// Set marshals the provided data into JSON and sets it in Redis with the specified key and expiration duration.
func (c Client) Set(ctx context.Context, key string, duration time.Duration, data interface{}) error {
	// Marshal the data into JSON.
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Set the JSON-encoded data in Redis with the specified key and expiration duration.
	return c.redisClient.Set(ctx, key, dataBytes, duration).Err()
}
