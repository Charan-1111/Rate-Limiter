package logic

import (
	"context"
	"goapp/algorithms"
	"goapp/models"
	"goapp/services"
	"goapp/store"
	"goapp/utils"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

func GetLimiter(ctx context.Context, db *store.Db, rdb *redis.Client, config *utils.Config, log zerolog.Logger, limiterFactory algorithms.LimiterFactory, cache *services.Cache, cb *services.CircuitBreaker, scope, identifier, rateLimitType string) (*models.LimiterResponse, error) {
	limiter, err := limiterFactory.GetLimiter(ctx, db, log, scope, identifier, rateLimitType, config.Queries.Fetch.FetchPolicyByKey, cache)
	if err != nil {
		log.Error().Err(err).Msg("Error getting the limiter interface")
		return nil, err
	}

	return limiter.Allow(ctx, rdb, cb, log, scope, identifier)
}
