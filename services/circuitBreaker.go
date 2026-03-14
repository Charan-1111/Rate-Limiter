package services

import (
	"goapp/constants"

	"github.com/sony/gobreaker"
)

type CircuitBreaker struct {
	Cb *gobreaker.CircuitBreaker
}

func NewCircuitBreaker() *CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        "CircuitBreaker",
		MaxRequests: 2,                                // half-open state allowed requests
		Interval:    constants.CircuitBreakerInterval, // time to reset counts
		Timeout:     constants.CircuitBreakerTimeout,  // time to stay in open state

		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= constants.ConsecutiveFailuresThreshold
		},
	}

	return &CircuitBreaker{
		Cb: gobreaker.NewCircuitBreaker(settings),
	}
}
