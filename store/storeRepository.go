package store

type LimiterRepository interface {
	// FetchPolicies() map[string]*services.PolicySchema
	// FetchPolicyByKey(key string) (*services.PolicySchema, bool)
}

// type StoreStruct struct {
// 	db  *pgxpool.Pool
// 	rdb *redis.Client
// 	log zerolog.Logger
// }
