package store

import (
	"context"
	"net"
	"time"

	"goapp/metrics"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type RedisConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type metricsHook struct{}

func (metricsHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

func (metricsHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now()
		err := next(ctx, cmd)
		duration := time.Since(start).Seconds()

		metrics.RedisLatency.Observe(duration)
		if err != nil && err != redis.Nil {
			metrics.RedisErrors.Inc()
		}
		return err
	}
}

func (metricsHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		start := time.Now()
		err := next(ctx, cmds)
		duration := time.Since(start).Seconds()

		metrics.RedisLatency.Observe(duration)
		if err != nil && err != redis.Nil {
			metrics.RedisErrors.Inc()
		}
		return err
	}
}

func InitRedis(redisDetails *RedisConfig, log zerolog.Logger) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisDetails.Host + ":" + redisDetails.Port,
	})

	rdb.AddHook(metricsHook{})

	return rdb
}
