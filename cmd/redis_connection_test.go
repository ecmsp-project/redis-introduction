package main

import (
	"context"
	"fmt"
	"testing"
)

func TestRedisServerConnection(t *testing.T) {
	rdb := setupTestRedisClient()
	defer rdb.Close()

	ctx := context.Background()

	err := rdb.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
		t.Fatal("error - unable to set key and value in redis database")
	}

	value, err := rdb.Get(ctx, "foo").Result()
	if err != nil {
		t.Fatal("error - unable to fetch value of given key in redis database")
	}

	fmt.Println("foo", value)
}
