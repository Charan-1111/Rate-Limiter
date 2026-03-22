package services

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func ExecuteLuaScript(ctx context.Context, rdb *redis.Client, keys []string, policy *PolicySchema) {

}