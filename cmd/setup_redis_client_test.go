package main

import "github.com/redis/go-redis/v9"

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
