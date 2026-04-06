package services

import (
	"context"
	"goapp/constants"
	"goapp/store"

	"github.com/dgraph-io/ristretto"
	"github.com/rs/zerolog"
)

type PolicySchema struct {
	Scope      string `json:"scope"`
	Identifier string `json:"identifier"`
	Limit      int    `json:"limit"`
	Window     string `json:"window"`
	Burst      int    `json:"burst"`
	Algorithm  string `json:"algorithm"`
}

type Cache struct {
	data *ristretto.Cache
}

func NewCache() *Cache {
	c, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e6,
		MaxCost:     1 << 28,
		BufferItems: 64,
	})
	return &Cache{
		data: c,
	}
}

func (c *Cache) LoadCache(ctx context.Context, log zerolog.Logger, db *store.Db, query string) {
	policies := FetchPolicies(ctx, db, log, query)

	for policyKey, policy := range policies {
		c.data.SetWithTTL(policyKey, policy, 1, constants.PolicyCacheDuration)
	}
}

func (c *Cache) GetPolicy(ctx context.Context, db *store.Db, log zerolog.Logger, scope, identifier, query string) (*PolicySchema, bool) {
	cacheKey := scope + ":" + identifier

	if val, found := c.data.Get(cacheKey); found {
		if cachedPolicy, ok := val.(*PolicySchema); ok {
			return cachedPolicy, true
		}
	}

	dbPolicy, exists := FetchPolicyByKey(ctx, db, log, query, cacheKey)
	if exists && dbPolicy != nil {
		c.data.SetWithTTL(cacheKey, dbPolicy, 1, constants.PolicyCacheDuration)
		return dbPolicy, true
	}

	return nil, false
}
