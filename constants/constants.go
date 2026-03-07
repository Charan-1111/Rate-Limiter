package constants

const (
	// ratelimit keys
	KeyRateLimit = "rate_limit"
	KeyRateLimitType = "type"
	KeyAlgo = "algo"

	// Values
	ValeTypeMemory = "memory"
	ValueTypeRedis = "redis"
	

	// Algorithms
	AlgorithmTokenBucket = "token_bucket"
	AlgorithmLeakyBucket = "leaky_bucket"
	AlgorithmFixedWindow = "fixed_window"
	AlgorithmSlidingWindow = "sliding_window"
)