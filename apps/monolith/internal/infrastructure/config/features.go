package config

import (
	"os"
	"strings"
)

// FeatureFlags controls optional module behavior.
type FeatureFlags struct {
	UseEventBus bool
}

// LoadFeatureFlags reads feature flag values from environment variables.
func LoadFeatureFlags() FeatureFlags {
	return FeatureFlags{
		UseEventBus: getEnvBool("USE_EVENT_BUS", false),
	}
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return strings.ToLower(value) == "true" || value == "1"
}
