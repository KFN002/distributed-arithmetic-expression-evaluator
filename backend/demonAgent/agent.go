package demonAgent

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

func RedisQueueHandler(redisClient *redis.Client) {
	ctx := context.Background()
	for {
		expression, err := redisClient.BRPop(ctx, 0, "expressions_queue").Result()
		if err != nil {
			log.Println("Error retrieving expression from Redis queue:", err)
			continue
		}
		fmt.Println(expression)

		if err != nil {
			log.Println("Error processing expression:", err)
			continue
		}
	}
}
