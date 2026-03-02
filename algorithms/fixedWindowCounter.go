package algorithms

import (
	"context"
	"fmt"
	"goapp/constants"
	"goapp/lua"
	"goapp/store"
	"time"

	"github.com/redis/go-redis/v9"
)

type FixedCounterStore struct {
	WindowIndex int64
	Allowed     int64
}

type FixedCounter struct {
	window   time.Duration
	capacity int64
}

func GetNewFixedWindowCounter(window time.Duration, capacity int64) *FixedCounter {
	return &FixedCounter{
		window:   window,
		capacity: capacity,
	}
}

func (fc *FixedCounter) Allow(ctx context.Context, tenandId, userId string) (bool, error) {
	now := time.Now().UnixNano()

	redisKey := fmt.Sprintf("%s:%s:%s:%s", constants.KeyRateLimit, constants.AlgorithmFixedWindow, tenandId, userId)

	fwcScript := redis.NewScript(lua.GetFixedWindowCounterScript())

	_, err := fwcScript.Run(ctx, store.Rdb, []string{redisKey}, fc.capacity, fc.window, now, 1).Result()
	if err != nil {
		fmt.Println("Error calling the fixed window counter script, rejecting the request : ", err)
		return false, err
	} else {
		fmt.Println("Accepting the request")
	}
	return true, nil
}
