package algorithms

import (
	"context"
	"errors"
	"fmt"
	"goapp/constants"
	"goapp/services"
	"goapp/utils"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type TokenBucketStore struct {
	tokens   float64
	lastFill time.Time
}

type TokenBucket struct {
	capacity float64
	fillRate float64
	tokens   map[string]*TokenBucketStore
	mu       sync.Mutex
}

func NewTokenBucketMem(capacity, fillRate float64) *TokenBucket {
	return &TokenBucket{
		capacity: capacity,
		fillRate: fillRate,
		tokens:   make(map[string]*TokenBucketStore),
		mu:       sync.Mutex{},
	}
}

func (tb *TokenBucket) Allow(ctx context.Context, rdb *redis.Client, cb *services.CircuitBreaker, log zerolog.Logger, scope, identifier string) (bool, error) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// make the key
	key := utils.StringBuilder(constants.KeyRateLimit, constants.AlgorithmTokenBucket, scope, identifier)
	now := time.Now()

	// Fetch from the cache
	tokenStore, ok := tb.tokens[key]
	if !ok {
		// create a new store if this key hasn't been seen before
		fmt.Println("Not found in the cache")
		tokenStore = &TokenBucketStore{
			tokens:   tb.capacity,
			lastFill: now,
		}
	}

	// fill the tokens
	tokenStore.tokens = tokenStore.tokens + (now.Sub(tokenStore.lastFill).Seconds() * tb.fillRate)
	if tokenStore.tokens >= tb.capacity {
		tokenStore.tokens = tb.capacity
	}

	tokenStore.lastFill = now

	// check if the bucket is empty
	if tokenStore.tokens == 0 {
		fmt.Println("Request is getting rejected, bucket is empty")
		return false, errors.New("Request is getting rejected")
	}

	tokenStore.tokens -= 1

	// store this in the cache
	tb.tokens[key] = tokenStore

	return true, nil
}
