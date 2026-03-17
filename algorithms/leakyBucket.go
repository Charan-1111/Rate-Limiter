package algorithms

import (
	"context"
	"goapp/constants"
	"goapp/lua"
	"goapp/services"
	"goapp/utils"
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

func (lb *LeakyBucketRedis) Allow(ctx context.Context, rdb *redis.Client, cb *services.CircuitBreaker, log zerolog.Logger, scope, identifier string) (bool, error) {
	// read data from redis
	redisKey := utils.StringBuilder(constants.KeyRateLimit, constants.AlgorithmLeakyBucket, scope, identifier)

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
