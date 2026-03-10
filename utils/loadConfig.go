package utils

import (
	"goapp/models"
	"goapp/store"
	"os"

	"github.com/bytedance/sonic"
)

type Config struct {
	Ports      models.Ports      `json:"ports"`
	Database   store.Database    `json:"database"`
	Redis      store.RedisConfig `json:"redis"`
	Tables     map[string]string `json:"tables"`
	MaxTokens  float64           `json:"maxTokens"`
	RefillRate float64           `json:"refillRate"`
}

func (config *Config) LoadConfig(filePath string) error {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := sonic.Unmarshal(fileBytes, config); err != nil {
		return err
	}

	return nil
}
