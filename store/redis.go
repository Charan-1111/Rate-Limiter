package store

import (
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type RedisConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

func InitRedis(redisDetails *RedisConfig, log zerolog.Logger) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisDetails.Host + ":" + redisDetails.Port,
	})

	return rdb
}
