package store

import (
	"goapp/services"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type LimiterRepository interface {
	FetchPolicies() map[string]*services.PolicySchema
	FetchPolicyByKey(key string) (*services.PolicySchema, bool)
}

type StoreStruct struct {
	db  *pgxpool.Pool
	rdb *redis.Client
	log zerolog.Logger
}
