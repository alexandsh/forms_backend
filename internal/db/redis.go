package db

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

func NewRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
}

func PingRedis(rdb *redis.Client) error {
	return rdb.Ping(context.Background()).Err()
}
