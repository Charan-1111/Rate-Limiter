package logic

import (
	"context"
	"fmt"
	"goapp/algorithms"

	"github.com/redis/go-redis/v9"
)

func GetLimiter(rdb *redis.Client, limiterFactory algorithms.LimiterFactory, limiterType, algorithm string) {
	limiter, err := limiterFactory.GetLimiter(limiterType, algorithm)
	if err != nil {
		fmt.Println("Error getting the limiter : ", err)
	}

	limiter.Allow(context.Background(), rdb, "tenant1", "user1")
}
