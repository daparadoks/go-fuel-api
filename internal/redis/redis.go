package redisService

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func Get(ctx context.Context, key string) string {
	_rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       1,  // use default DB
	})
	val, err := _rdb.Get(ctx, key).Result()
	if err != nil {
		return ""
	}

	return val
}

func Set(ctx context.Context, key string, val string) {
	if len(val) == 0 {
		return
	}
	_rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       1,  // use default DB
	})
	err := _rdb.Set(ctx, key, val, 0).Err()
	if err != nil {
		panic(err)
	}
}

func Delete(ctx context.Context, key string) {
	_rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       1,  // use default DB
	})
	err := _rdb.Del(ctx, key).Err()
	if err != nil {
		panic(err)
	}
}
