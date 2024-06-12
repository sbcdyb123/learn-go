package ioc

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func InitRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	for err := redisClient.Ping(context.Background()).Err(); err != nil; {
		panic(err)
	}
	return redisClient
}
