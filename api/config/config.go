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
	App    App
	OpenAI OpenAI
	Redis  RedisClient
}

// App app configuration
type App struct {
	ChatLimitDuration int64  `env:"CHAT_LIMIT_DURATION" envDefault:"3"`
	Port              string `env:"PORT" envDefault:":5000"`
}

type OpenAI struct {
	ChatCompletionAPIURL string `env:"OPENAI_CHAT_COMPLETION_URL" envDefault:"https://api.openai.com/v1/chat/completions"`
	ChatCompletionModel  string `env:"OPENAI_CHAT_COMPLETION_MODEL" envDefault:"gpt-3.5-turbo-1106"`

	APIKey string `env:"OPENAI_API_KEY" envDefault:"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"`
}

// RedisClient represents the redis connection
type RedisClient struct {
	RedisAddr    string `env:"REDIS_CLIENT_ADDR" envDefault:"localhost:6079"`
	MinIdleConns int    `env:"REDIS_CLIENT_MIN_IDLE_CONNS" envDefault:"10"`
	PoolSize     int    `env:"REDIS_CLIENT_POOL_SIZE" envDefault:"10"`
	PoolTimeout  int    `env:"REDIS_CLIENT_POOL_TIMEOUT" envDefault:"3"`
}
