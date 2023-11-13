package database

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisCtx = context.Background()

func RedisClient(dbName int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASS"),
		DB: dbName,
	})
}