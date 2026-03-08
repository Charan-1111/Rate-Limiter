package algorithms

import (
	"context"
	"fmt"
	"goapp/constants"
	"goapp/lua"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type SlidingWindowStoreRedis struct {
	CurrentCnt  int
	PreviousCnt int
	WindowStart time.Time
}
type SlidingWindowCounterRedis struct {
	window   time.Duration
	capacity int
	mu       sync.Mutex
}

func NewSlidingWindowCounter(window time.Duration, capacity int) *SlidingWindowCounterRedis {
	return &SlidingWindowCounterRedis{
		window:   window,
		capacity: capacity,
	}
}

func (sc *SlidingWindowCounterRedis) Allow(ctx context.Context, rdb *redis.Client, tenantId, userId string) (bool, error) {
	now := time.Now().UnixMilli()

	windowMs := sc.window.Milliseconds()
	// read data from the redis
	redisKey := fmt.Sprintf("%s:%s:%s:%s", constants.KeyRateLimit, constants.AlgorithmSlidingWindow, tenantId, userId)

	swcScript := redis.NewScript(lua.GetSlidingWindowScript())

	_, err := swcScript.Run(ctx, rdb, []string{redisKey}, sc.capacity, windowMs, now, 1).Result()
	if err != nil {
		fmt.Println("Error calling the sliding window counter script, rejecting the request : ", err)
		return false, err
	} else {
		fmt.Println("Accepting the request")
	}
	return true, nil
}
