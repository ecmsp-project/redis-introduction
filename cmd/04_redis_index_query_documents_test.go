package main

/*
	Redis allows to perform search on both JSON and hash objects. 
	Create an index and define a schema. Next, add data (JSONSet or HSet function depending on type) and create the query (allows to group, filter, count etc.)
*/

import (
	"context"
	"fmt"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestRedisIndexQueryDocuments(t *testing.T) {
	rdb := setupTestRedisClient()
	defer rdb.Close()
	ctx := context.Background()

	err := flushRedisDB(ctx, rdb)
	require.NoError(t, err)


	user1 := map[string]interface{}{
		"name":  "Paul John",
		"email": "paul.john@example.com",
		"age":   42,
		"city":  "London",
	}

	user2 := map[string]interface{}{
		"name":  "Eden Zamir",
		"email": "eden.zamir@example.com",
		"age":   29,
		"city":  "Tel Aviv",
	}

	user3 := map[string]interface{}{
		"name":  "Paul Zamir",
		"email": "paul.zamir@example.com",
		"age":   35,
		"city":  "Tel Aviv",
	}

	_, err = rdb.FTCreate(
		ctx,
		"idx:users",
		&redis.FTCreateOptions{
			OnJSON: true,
			Prefix: []interface{}{"user:"},
		},
		&redis.FieldSchema{
			FieldName: "$.name",
			As:        "name",
			FieldType: redis.SearchFieldTypeText,
		},
		&redis.FieldSchema{
			FieldName: "$.city",
			As:        "city",
			FieldType: redis.SearchFieldTypeTag,
		},
		&redis.FieldSchema{
			FieldName: "$.age",
			As:        "age",
			FieldType: redis.SearchFieldTypeNumeric,
		},
	).Result()
	require.NoError(t, err)

	_, err = rdb.JSONSet(ctx, "user:1", "$", user1).Result()
	require.NoError(t, err)

	_, err = rdb.JSONSet(ctx, "user:2", "$", user2).Result()
	require.NoError(t, err)

	_, err = rdb.JSONSet(ctx, "user:3", "$", user3).Result()
	require.NoError(t, err)

	findPaulResult, err := rdb.FTSearch(
		ctx,
		"idx:users",
		"Paul @age:[30 40]",
	).Result()

	require.NoError(t, err)
	require.NotEmpty(t, findPaulResult)

	fmt.Println(findPaulResult)
}
