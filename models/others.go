package models

type Ports struct {
	FiberServer string `json:"fiberServer"`
}

type LimiterResponse struct {
	Allowed       bool  `json:"allowed"`
	RetryAfter    int64 `json:"retryAfter"`
	CurrentTokens int64 `json:"remaining"`
}
