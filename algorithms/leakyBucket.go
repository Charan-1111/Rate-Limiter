package algorithms

import (
	"context"
	"fmt"
	"goapp/constants"
	"goapp/lua"
	"goapp/services"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type LeakyTokensRedis struct {
	Tokens   float64
	LastLeak time.Time
}

type LeakyBucketRedis struct {
	MaxTokens float64
	LeakRate  float64
	mu        sync.Mutex
}

func NewLeakyBucket(maxTokens, leakRate float64) *LeakyBucketRedis {
	return &LeakyBucketRedis{
		MaxTokens: maxTokens,
		LeakRate:  leakRate,
	}
}

func (lb *LeakyBucketRedis) Allow(ctx context.Context, rdb *redis.Client, cb *services.CircuitBreaker, log zerolog.Logger, tenantId, userId string) (bool, error) {
	// read data from redis
	redisKey := fmt.Sprintf("%s:%s:%s:%s", constants.KeyRateLimit, constants.AlgorithmLeakyBucket, tenantId, userId)

	leakyScript := redis.NewScript(lua.GetLeakyBucketScript())

	now := float64(time.Now().UnixNano()) / 1e9

	// _, err := leakyScript.Run(ctx, rdb, []string{redisKey}, lb.MaxTokens, lb.LeakRate, now, 1).Result()
	// if err != nil {
	// 	fmt.Println("Error running the script, rejecting the request : ", err)
	// 	return false, err
	// } else {
	// 	fmt.Println("request accepted")
	// }

	_, err := cb.Cb.Execute(func() (any, error) {
		return leakyScript.Run(ctx, rdb, []string{redisKey}, lb.MaxTokens, lb.LeakRate, now, 1).Result()
	})

	if err != nil {
		log.Error().Err(err).Msg("Error running the script, rejecting the request")
		return false, err
	} else {
		log.Info().Msg("Accepting the request")
	}

	return true, nil
}
