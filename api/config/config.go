package config

import (
	"github.com/caarlos0/env/v6"
)

// Instance the instance
var Instance *Config

// Initial new the impl config obj
func Initial() error {
	Instance = &Config{}

	if err := env.Parse(Instance); err != nil {
		return err
	}

	return nil
}

// Config represents a config instance
type Config struct {
	App   App
	Redis RedisClient
}

// App app configuration
type App struct {
	ChatLimitDuration int64  `env:"CHAT_LIMIT_DURATION" envDefault:"3"`
	Port              string `env:"PORT" envDefault:":5000"`
}

// RedisClient represents the redis connection
type RedisClient struct {
	RedisAddr    string `env:"REDIS_CLIENT_ADDR" envDefault:"localhost:6079"`
	MinIdleConns int    `env:"REDIS_CLIENT_MIN_IDLE_CONNS" envDefault:"10"`
	PoolSize     int    `env:"REDIS_CLIENT_POOL_SIZE" envDefault:"10"`
	PoolTimeout  int    `env:"REDIS_CLIENT_POOL_TIMEOUT" envDefault:"3"`
}
