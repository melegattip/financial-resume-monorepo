package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func clearEnv() {
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_SSLMODE")
	os.Unsetenv("DB_MAX_IDLE_CONNS")
	os.Unsetenv("DB_MAX_OPEN_CONNS")
	os.Unsetenv("PORT")
	os.Unsetenv("APP_ENV")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("CORS_ALLOW_ORIGINS")
}

func TestLoad_WithDatabaseURL(t *testing.T) {
	clearEnv()
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/testdb?sslmode=disable")

	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, "postgres://user:pass@localhost:5432/testdb?sslmode=disable", cfg.DatabaseURL)
	assert.Equal(t, "postgres://user:pass@localhost:5432/testdb?sslmode=disable", cfg.DSN())
}

func TestLoad_WithIndividualDBVars(t *testing.T) {
	clearEnv()
	t.Setenv("DB_USER", "testuser")
	t.Setenv("DB_PASSWORD", "testpass")
	t.Setenv("DB_HOST", "dbhost")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_NAME", "mydb")
	t.Setenv("DB_SSLMODE", "require")

	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, "testuser", cfg.DBUser)
	assert.Equal(t, "testpass", cfg.DBPassword)
	assert.Contains(t, cfg.DSN(), "host=dbhost")
	assert.Contains(t, cfg.DSN(), "port=5433")
	assert.Contains(t, cfg.DSN(), "dbname=mydb")
	assert.Contains(t, cfg.DSN(), "sslmode=require")
}

func TestLoad_DatabaseURLTakesPrecedence(t *testing.T) {
	clearEnv()
	t.Setenv("DATABASE_URL", "postgres://url@host/db")
	t.Setenv("DB_USER", "ignored")
	t.Setenv("DB_PASSWORD", "ignored")

	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, "postgres://url@host/db", cfg.DSN())
}

func TestLoad_MissingRequiredVars(t *testing.T) {
	clearEnv()

	_, err := Load()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "DB_USER")
}

func TestLoad_MissingPassword(t *testing.T) {
	clearEnv()
	t.Setenv("DB_USER", "testuser")

	_, err := Load()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "DB_PASSWORD")
}

func TestLoad_Defaults(t *testing.T) {
	clearEnv()
	t.Setenv("DATABASE_URL", "postgres://u:p@h/d")

	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, "8080", cfg.ServerPort)
	assert.Equal(t, "development", cfg.Environment)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, "http://localhost:3000", cfg.CORSAllowOrigins)
	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, "5432", cfg.DBPort)
	assert.Equal(t, "financial_resume", cfg.DBName)
	assert.Equal(t, "disable", cfg.DBSSLMode)
	assert.Equal(t, 10, cfg.DBMaxIdleConns)
	assert.Equal(t, 25, cfg.DBMaxOpenConns)
}

func TestLoad_CustomValues(t *testing.T) {
	clearEnv()
	t.Setenv("DATABASE_URL", "postgres://u:p@h/d")
	t.Setenv("PORT", "9090")
	t.Setenv("APP_ENV", "production")
	t.Setenv("LOG_LEVEL", "error")
	t.Setenv("DB_MAX_IDLE_CONNS", "5")
	t.Setenv("DB_MAX_OPEN_CONNS", "50")

	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, "9090", cfg.ServerPort)
	assert.Equal(t, "production", cfg.Environment)
	assert.Equal(t, "error", cfg.LogLevel)
	assert.Equal(t, 5, cfg.DBMaxIdleConns)
	assert.Equal(t, 50, cfg.DBMaxOpenConns)
}
