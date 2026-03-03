package memoryalgorithms

import (
	"context"
	"fmt"
	"goapp/constants"
	"sync"
	"time"
)

type FixedWindowStore struct {
	windowIndex int
	tokens      int
}

type FixedWindow struct {
	capacity int
	window   time.Time
	tokens   map[string]*FixedWindowStore
	mu       sync.Mutex
}

func (fw *FixedWindow) Allow(ctx context.Context, tenantId, userId string) (bool, error) {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	key := fmt.Sprintf("%s:%s:%s:%s", constants.KeyRateLimit, constants.AlgorithmFixedWindow, tenantId, userId)
	now := time.Now()

	// fetch the data from the cache

	tokenStore, ok := fw.tokens[key]
	if !ok {
		
	}

	return true, nil
}
