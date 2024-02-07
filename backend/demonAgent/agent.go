package demonAgent

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

func RedisQueueHandler(redisClient *redis.Client) {
	ctx := context.Background()

	// Infinite loop to continuously process expressions from the queue
	for {
		// Retrieve expression from the Redis queue
		expression, err := redisClient.BRPop(ctx, 0, "expressions_queue").Result()
		if err != nil {
			log.Println("Error retrieving expression from Redis queue:", err)
			// Handle error appropriately (e.g., retry mechanism, logging)
			continue
		}

		// Process the retrieved expression
		fmt.Println(expression)

		if err != nil {
			log.Println("Error processing expression:", err)
			// Handle error appropriately (e.g., retry mechanism, logging)
			continue
		}
	}
}
