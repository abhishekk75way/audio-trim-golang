package config

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

func ConnectRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})

	if err := Redis.Ping(context.Background()).Err(); err != nil {
		panic("Failed to connect Redis")
	}
}
