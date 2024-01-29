package redis

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/tanmaij/zylo/config"
)

// NewRedisClient new redis client instance
func NewRedisClient() (*redis.Client, error) {
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
		return nil, err
	}

	log.Println("Connected to Redis:", pong)
	return rdb, nil
}
