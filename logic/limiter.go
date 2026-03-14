package logic

import (
	"context"
	"goapp/algorithms"
	"goapp/services"
	"goapp/utils"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

func GetLimiter(ctx context.Context, db *pgxpool.Pool, rdb *redis.Client, config *utils.Config, log zerolog.Logger, limiterFactory algorithms.LimiterFactory, cache *services.Cache, scope, identifier, rateLimitType string) {
	limiter, err := limiterFactory.GetLimiter(ctx, db, log, scope, identifier, rateLimitType, config.Queries.Fetch.FetchPolicyByKey, cache)
	if err != nil {
		log.Error().Err(err).Msg("Error getting the limiter interface")
	}

	limiter.Allow(ctx, rdb, scope, identifier)
}
