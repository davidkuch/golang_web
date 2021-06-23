package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// docker run -p 6379:6379 redis

var ctx = context.Background()

func redisSetSession(name string, uuid string) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, uuid, name, 6000*5).Err()
	if err != nil {
		panic(err)
	}
}

func redisGetSession(uuid string) string {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	name, err := rdb.Get(ctx, uuid).Result()
	if err != nil {
		fmt.Println(err)
	}
	return name

}
