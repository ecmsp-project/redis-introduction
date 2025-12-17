package main

/*

 */

import (
	"context"
	"fmt"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestRedisPipelinesTest(t *testing.T) {
	rdb := setupTestRedisClient()
	defer rdb.Close()
	ctx := context.Background()

	err := flushRedisDB(ctx, rdb)
	require.NoError(t, err)

	pipe := rdb.Pipeline()

	for i := 0; i < 5; i++ {
		pipe.Set(ctx, fmt.Sprintf("seat:%v", i), fmt.Sprintf("#%v", i), 0)
	}

	cmds, err := pipe.Exec(ctx)
	require.NoError(t, err)

	for _, c := range cmds {
		fmt.Printf("%v;\n", c.(*redis.StatusCmd).Val())
	}

	pipe = rdb.Pipeline()

	get0Result := pipe.Get(ctx, "seat:0")
	get3Result := pipe.Get(ctx, "seat:3")
	get4Result := pipe.Get(ctx, "seat:4")

	_, err = pipe.Exec(ctx)
	require.NoError(t, err)

	fmt.Println(get0Result.Val())
	fmt.Println(get3Result.Val())
	fmt.Println(get4Result.Val())
}
