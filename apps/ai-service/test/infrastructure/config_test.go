package infrastructure

import (
	"os"
	"testing"

	"github.com/financial-ai-service/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ConfigTestSuite suite de tests para la configuración
type ConfigTestSuite struct {
	suite.Suite
	originalEnv map[string]string
}

// SetupSuite guarda las variables de entorno originales
func (suite *ConfigTestSuite) SetupSuite() {
	suite.originalEnv = make(map[string]string)
	envVars := []string{
		"PORT", "HOST", "OPENAI_API_KEY", "USE_AI_MOCK", "REDIS_URL",
		"REDIS_PASSWORD", "REDIS_DB", "CACHE_DEFAULT_TTL_MINUTES", "CACHE_INSIGHTS_TTL_HOURS",
	}

	for _, envVar := range envVars {
		if value, exists := os.LookupEnv(envVar); exists {
			suite.originalEnv[envVar] = value
		}
	}
}

// TearDownSuite restaura las variables de entorno originales
func (suite *ConfigTestSuite) TearDownSuite() {
	envVars := []string{
		"PORT", "HOST", "OPENAI_API_KEY", "USE_AI_MOCK", "REDIS_URL",
		"REDIS_PASSWORD", "REDIS_DB", "CACHE_DEFAULT_TTL_MINUTES", "CACHE_INSIGHTS_TTL_HOURS",
	}

	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}

	for key, value := range suite.originalEnv {
		os.Setenv(key, value)
	}
}

// SetupTest limpia las variables de entorno para cada test
func (suite *ConfigTestSuite) SetupTest() {
	envVars := []string{
		"PORT", "HOST", "OPENAI_API_KEY", "USE_AI_MOCK", "REDIS_URL",
		"REDIS_PASSWORD", "REDIS_DB", "CACHE_DEFAULT_TTL_MINUTES", "CACHE_INSIGHTS_TTL_HOURS",
	}

	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}

// TestLoad_DefaultValues testa la carga con valores por defecto
func (suite *ConfigTestSuite) TestLoad_DefaultValues() {
	// Act
	cfg, err := config.Load()

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), cfg)
	assert.Equal(suite.T(), "8082", cfg.Server.Port)
	assert.Equal(suite.T(), "localhost", cfg.Server.Host)
	assert.Equal(suite.T(), "", cfg.OpenAI.APIKey)
	assert.True(suite.T(), cfg.OpenAI.UseMock)
	assert.Equal(suite.T(), "redis://localhost:6379", cfg.Redis.URL)
	assert.Equal(suite.T(), "", cfg.Redis.Password)
	assert.Equal(suite.T(), 0, cfg.Redis.DB)
	assert.Equal(suite.T(), 30, cfg.Cache.DefaultTTLMinutes)
	assert.Equal(suite.T(), 20, cfg.Cache.InsightsTTLHours)
}

// TestLoad_WithEnvironmentVariables testa la carga con variables de entorno
func (suite *ConfigTestSuite) TestLoad_WithEnvironmentVariables() {
	// Arrange
	envVars := map[string]string{
		"PORT":                      "3000",
		"HOST":                      "0.0.0.0",
		"OPENAI_API_KEY":            "sk-test123456789",
		"USE_AI_MOCK":               "false",
		"REDIS_URL":                 "redis://custom:6379",
		"REDIS_PASSWORD":            "secret",
		"REDIS_DB":                  "1",
		"CACHE_DEFAULT_TTL_MINUTES": "60",
		"CACHE_INSIGHTS_TTL_HOURS":  "48",
	}

	for key, value := range envVars {
		os.Setenv(key, value)
	}

	// Act
	cfg, err := config.Load()

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), cfg)
	assert.Equal(suite.T(), "3000", cfg.Server.Port)
	assert.Equal(suite.T(), "0.0.0.0", cfg.Server.Host)
	assert.Equal(suite.T(), "sk-test123456789", cfg.OpenAI.APIKey)
	assert.False(suite.T(), cfg.OpenAI.UseMock)
	assert.Equal(suite.T(), "redis://custom:6379", cfg.Redis.URL)
	assert.Equal(suite.T(), "secret", cfg.Redis.Password)
	assert.Equal(suite.T(), 1, cfg.Redis.DB)
	assert.Equal(suite.T(), 60, cfg.Cache.DefaultTTLMinutes)
	assert.Equal(suite.T(), 48, cfg.Cache.InsightsTTLHours)
}

// TestLoad_InvalidValues testa el manejo de valores inválidos
func (suite *ConfigTestSuite) TestLoad_InvalidValues() {
	testCases := []struct {
		name     string
		envVar   string
		value    string
		expected interface{}
	}{
		{
			name:     "USE_AI_MOCK inválido",
			envVar:   "USE_AI_MOCK",
			value:    "invalid",
			expected: true, // Debería usar el default
		},
		{
			name:     "REDIS_DB inválido",
			envVar:   "REDIS_DB",
			value:    "invalid",
			expected: 0, // Debería usar el default
		},
		{
			name:     "CACHE_DEFAULT_TTL_MINUTES inválido",
			envVar:   "CACHE_DEFAULT_TTL_MINUTES",
			value:    "invalid",
			expected: 30, // Debería usar el default
		},
		{
			name:     "CACHE_INSIGHTS_TTL_HOURS inválido",
			envVar:   "CACHE_INSIGHTS_TTL_HOURS",
			value:    "invalid",
			expected: 20, // Debería usar el default
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Arrange
			os.Setenv(tc.envVar, tc.value)

			// Act
			cfg, err := config.Load()

			// Assert
			assert.NoError(suite.T(), err)
			assert.NotNil(suite.T(), cfg)

			switch tc.envVar {
			case "USE_AI_MOCK":
				assert.Equal(suite.T(), tc.expected, cfg.OpenAI.UseMock)
			case "REDIS_DB":
				assert.Equal(suite.T(), tc.expected, cfg.Redis.DB)
			case "CACHE_DEFAULT_TTL_MINUTES":
				assert.Equal(suite.T(), tc.expected, cfg.Cache.DefaultTTLMinutes)
			case "CACHE_INSIGHTS_TTL_HOURS":
				assert.Equal(suite.T(), tc.expected, cfg.Cache.InsightsTTLHours)
			}

			// Cleanup
			os.Unsetenv(tc.envVar)
		})
	}
}

// TestLoad_EmptyValues testa el manejo de valores vacíos
func (suite *ConfigTestSuite) TestLoad_EmptyValues() {
	// Arrange
	envVars := map[string]string{
		"PORT":           "",
		"HOST":           "",
		"OPENAI_API_KEY": "",
		"REDIS_URL":      "",
		"REDIS_PASSWORD": "",
	}

	for key, value := range envVars {
		os.Setenv(key, value)
	}

	// Act
	cfg, err := config.Load()

	// Assert - debería usar valores por defecto cuando están vacíos
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), cfg)
	assert.Equal(suite.T(), "8082", cfg.Server.Port)
	assert.Equal(suite.T(), "localhost", cfg.Server.Host)
	assert.Equal(suite.T(), "", cfg.OpenAI.APIKey)
	assert.Equal(suite.T(), "redis://localhost:6379", cfg.Redis.URL)
	assert.Equal(suite.T(), "", cfg.Redis.Password)
}

// TestLoad_BooleanValues testa diferentes valores booleanos
func (suite *ConfigTestSuite) TestLoad_BooleanValues() {
	testCases := []struct {
		name     string
		value    string
		expected bool
	}{
		{
			name:     "true string",
			value:    "true",
			expected: true,
		},
		{
			name:     "false string",
			value:    "false",
			expected: false,
		},
		{
			name:     "1 como true",
			value:    "1",
			expected: true,
		},
		{
			name:     "0 como false",
			value:    "0",
			expected: false,
		},
		{
			name:     "TRUE en mayúsculas",
			value:    "TRUE",
			expected: true,
		},
		{
			name:     "FALSE en mayúsculas",
			value:    "FALSE",
			expected: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Arrange
			os.Setenv("USE_AI_MOCK", tc.value)

			// Cuando mock está deshabilitado, necesitamos proporcionar API key
			if !tc.expected {
				os.Setenv("OPENAI_API_KEY", "sk-test-key-for-testing")
			}

			// Act
			cfg, err := config.Load()

			// Assert
			assert.NoError(suite.T(), err)
			assert.NotNil(suite.T(), cfg)
			assert.Equal(suite.T(), tc.expected, cfg.OpenAI.UseMock)

			// Cleanup
			os.Unsetenv("USE_AI_MOCK")
			if !tc.expected {
				os.Unsetenv("OPENAI_API_KEY")
			}
		})
	}
}

// TestLoad_NumericValues testa diferentes valores numéricos
func (suite *ConfigTestSuite) TestLoad_NumericValues() {
	testCases := []struct {
		name     string
		envVar   string
		value    string
		expected int
	}{
		{
			name:     "REDIS_DB válido",
			envVar:   "REDIS_DB",
			value:    "5",
			expected: 5,
		},
		{
			name:     "TTL minutos",
			envVar:   "CACHE_DEFAULT_TTL_MINUTES",
			value:    "120",
			expected: 120,
		},
		{
			name:     "TTL horas",
			envVar:   "CACHE_INSIGHTS_TTL_HOURS",
			value:    "72",
			expected: 72,
		},
		{
			name:     "Valor cero",
			envVar:   "REDIS_DB",
			value:    "0",
			expected: 0,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Arrange
			os.Setenv(tc.envVar, tc.value)

			// Act
			cfg, err := config.Load()

			// Assert
			assert.NoError(suite.T(), err)
			assert.NotNil(suite.T(), cfg)

			switch tc.envVar {
			case "REDIS_DB":
				assert.Equal(suite.T(), tc.expected, cfg.Redis.DB)
			case "CACHE_DEFAULT_TTL_MINUTES":
				assert.Equal(suite.T(), tc.expected, cfg.Cache.DefaultTTLMinutes)
			case "CACHE_INSIGHTS_TTL_HOURS":
				assert.Equal(suite.T(), tc.expected, cfg.Cache.InsightsTTLHours)
			}

			// Cleanup
			os.Unsetenv(tc.envVar)
		})
	}
}

// TestLoad_ValidationErrors testa errores de validación
func (suite *ConfigTestSuite) TestLoad_ValidationErrors() {
	testCases := []struct {
		name    string
		envVars map[string]string
		error   string
	}{
		{
			name: "API key requerida sin mock",
			envVars: map[string]string{
				"USE_AI_MOCK":    "false",
				"OPENAI_API_KEY": "",
			},
			error: "OpenAI API key is required when not using mock",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Arrange
			for key, value := range tc.envVars {
				os.Setenv(key, value)
			}

			// Act
			cfg, err := config.Load()

			// Assert
			assert.Error(suite.T(), err)
			assert.Nil(suite.T(), cfg)
			assert.Contains(suite.T(), err.Error(), tc.error)

			// Cleanup
			for key := range tc.envVars {
				os.Unsetenv(key)
			}
		})
	}
}

// TestRunSuite ejecuta todos los tests del suite
func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

// Tests adicionales para casos específicos

// TestConfigValidation testa la validación básica de configuración
func TestConfigValidation(t *testing.T) {
	// Arrange & Act
	cfg, err := config.Load()

	// Assert - verificar que la configuración por defecto es válida
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.Server.Port)
	assert.NotEmpty(t, cfg.Server.Host)
	assert.NotEmpty(t, cfg.Redis.URL)
	assert.True(t, cfg.Cache.DefaultTTLMinutes > 0)
	assert.True(t, cfg.Cache.InsightsTTLHours > 0)
}

// TestConfigWithAPIKey testa configuración con API key real
func TestConfigWithAPIKey(t *testing.T) {
	// Arrange
	os.Setenv("OPENAI_API_KEY", "sk-real-api-key")
	os.Setenv("USE_AI_MOCK", "false")

	// Act
	cfg, err := config.Load()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "sk-real-api-key", cfg.OpenAI.APIKey)
	assert.False(t, cfg.OpenAI.UseMock)

	// Cleanup
	os.Unsetenv("OPENAI_API_KEY")
	os.Unsetenv("USE_AI_MOCK")
}
