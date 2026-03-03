package memoryalgorithms

import (
	"context"
	"errors"
	"fmt"
	"goapp/constants"
	"sync"
	"time"
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

type RateLimiter interface {
	Allow(ctx context.Context, tenantId string, userId string) (bool, error)
}

func (tb *TokenBucket) Allow(ctx context.Context, tenantId, userId string) (bool, error) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// make the key
	key := fmt.Sprintf("%s:%s:%s:%s", constants.KeyRateLimit, constants.AlgorithmTokenBucket, tenantId, userId)
	now := time.Now()

	// Fetch from the cache
	tokenStore, ok := tb.tokens[key]
	if !ok {
		fmt.Println("Not found in the cache")
		tokenStore.tokens = tb.capacity
		tokenStore.lastFill = now
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
