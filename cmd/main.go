package main

import (
	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := connectToRedisServer()

	// store := store.CreateRedisStore(rdb)

	defer rdb.Close()
}

func connectToRedisServer() *redis.Client {
	return redis.NewClient(
		&redis.Options{
			Addr:     "localhost:6379",
			Password: "admin",
			DB:       0,
			Protocol: 2,
		},
	)
}
