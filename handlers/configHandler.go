package handlers

import (
	"goapp/utils"

	"github.com/redis/go-redis/v9"
)

type ConfigHandler struct {
	config *utils.Config
	rdb    *redis.Client
}

func NewConfigHandler(config *utils.Config, rdb *redis.Client) *ConfigHandler {
	return &ConfigHandler{
		config: config,
		rdb:    rdb,
	}
}