package config

import (
	"github.com/melegattip/financial-resume-engine/pkg/config/types"
)

type Config struct{}

func (c *Config) GetStringSlice(key string, defaultValues []string) []string {
	return defaultValues
}

func (c *Config) GetInt(key string, defaultValue int) int {
	return defaultValue
}

func (c *Config) GetFloat64(key string, defaultValue float64) float64 {
	return defaultValue
}

func (c *Config) GetJSONPropertyAndUnmarshal(key string, structType interface{}) error {
	return nil
}

func Load() (types.Client, error) {
	return &Config{}, nil
}
