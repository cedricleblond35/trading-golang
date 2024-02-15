package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis() (*redis.Client, error) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "", // no password set
		DB:           0,  // use default DB
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})

	return rdb, rdb.Ping(ctx).Err()
}
