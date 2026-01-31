// Package configuration provides functionality for loading and managing configuration clients
package configuration

import (
	"context"
	"os"

	envconstants "github.com/melegattip/financial-resume-engine/internal/config/environment/constants"
	"github.com/melegattip/financial-resume-engine/internal/core/logs"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/logger"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/repository/configuration"
	"github.com/melegattip/financial-resume-engine/pkg/config"
	"github.com/melegattip/financial-resume-engine/pkg/config/configtest"
	"github.com/melegattip/financial-resume-engine/pkg/config/types"
)

// Client represents a configuration client that can load and access configuration values
type Client = types.Client

// LoadClient creates and returns a configuration client based on the current environment
// It returns a production client when in production environment, otherwise returns a test client
func LoadClient() Client {
	var configurationClient Client

	if os.Getenv("GO_ENVIRONMENT") == envconstants.Production {
		configurationClient, err := config.Load()
		if err != nil {
			logger.Error(context.Background(), err, logs.ErrorLoadingConfiguration.GetMessage(), logs.Tags{})
			panic(err)
		}

		return configurationClient
	}

	configurationClient = configtest.Load(configuration.DefaultConfigPath)

	return configurationClient
}
