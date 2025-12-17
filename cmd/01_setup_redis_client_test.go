package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func setupTestRedisClient() *redis.Client {
	return redis.NewClient(
		&redis.Options{
			Addr:     "localhost:6379",
			Password: "admin",
			DB:       0,
			Protocol: 2,
		},
	)
}

func flushRedisDB(ctx context.Context, rdb *redis.Client) error {
	return rdb.FlushDB(ctx).Err()
}
