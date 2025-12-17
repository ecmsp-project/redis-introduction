package main

/*
	Redis allows to store key-value pairs under a single key. It it ideal to store objects and recors with multiple fields, for example: records of different models, products, users. 
	Important functions are: HSet, HGet, HGetAll
*/

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

var hashFields map[string]string = map[string]string{
	"model": "Deimos",
	"brand": "Ergonom",
	"type":  "Enduro bikes",
	"price": "4972",
}

func TestRedisHashStorage(t *testing.T) {
	rdb := setupTestRedisClient()
	defer rdb.Close()
	ctx := context.Background()

	err := rdb.HSet(ctx, "bike:1", hashFields).Err()
	require.NoError(t, err)

	res1, err := rdb.HGet(ctx, "bike:1", "model").Result()
	require.NoError(t, err)
	require.NotEmpty(t, res1)
	require.Equal(t, "Deimos", res1)

	res2, err := rdb.HGetAll(ctx, "bike:1").Result()
	require.NoError(t, err)
	require.NotEmpty(t, res2)
}

func TestRedisHashStorageWithParsingToStruct(t *testing.T) {
	rdb := setupTestRedisClient()
	defer rdb.Close()
	ctx := context.Background()

	err := rdb.HSet(ctx, "bike:1", hashFields).Err()
	require.NoError(t, err)

	type BikeInfo struct {
		Model string `redis:"model"`
		UselessField string `redis:"useless_field"`
	}
	var res BikeInfo
	
	err = rdb.HGetAll(ctx, "bike:1").Scan(&res)
	require.NoError(t, err)
	require.NotEmpty(t, res)
	require.Equal(t, "Deimos", res.Model)
	require.Equal(t, "", res.UselessField)
}