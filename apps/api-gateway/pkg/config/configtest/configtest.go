package configtest

import (
	"github.com/melegattip/financial-resume-engine/pkg/config/types"
)

type TestConfig struct{}

func (c *TestConfig) GetStringSlice(key string, defaultValues []string) []string {
	return defaultValues
}

func (c *TestConfig) GetInt(key string, defaultValue int) int {
	return defaultValue
}

func (c *TestConfig) GetFloat64(key string, defaultValue float64) float64 {
	return defaultValue
}

func (c *TestConfig) GetJSONPropertyAndUnmarshal(key string, structType interface{}) error {
	return nil
}

func Load(configFile string) types.Client {
	return &TestConfig{}
}
