package algorithms

import (
	"context"
	"goapp/constants"
	"goapp/lua"
	"goapp/services"
	"goapp/utils"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type TokensRedis struct {
	Tokens     float64
	LastRefill time.Time
}

type TokenBucketRedis struct {
	MaxTokens  float64
	RefillRate float64
}

func NewTokenBucket(maxTokens, refillRate float64) *TokenBucketRedis {
	return &TokenBucketRedis{
		MaxTokens:  maxTokens,
		RefillRate: float64(refillRate),
	}
}

func (tb *TokenBucketRedis) Allow(ctx context.Context, rdb *redis.Client, cb *services.CircuitBreaker, log zerolog.Logger, scope, identifier string) (bool, error) {
	// tokens := &Tokens{}

	// get the information from the redis for the key
	redisKey := utils.StringBuilder(constants.KeyRateLimit, constants.AlgorithmTokenBucket, scope, identifier)

	tokenBucketScript := redis.NewScript(lua.GetTokenBucketScript())
	now := float64(time.Now().UnixNano()) / 1e9

	_, err := cb.Cb.Execute(func() (any, error) {
		return tokenBucketScript.Run(ctx, rdb, []string{redisKey}, tb.MaxTokens, tb.RefillRate, now, 1).Result()
	})

	if err != nil {
		log.Error().Err(err).Msg("Error running the token bucket script")
		return false, err
	} else {
		log.Info().Msg("Accepting the request")
	}

	return true, nil
}
