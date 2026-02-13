package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/http/handlers"
)

func TestNewRouter_HealthEndpoint(t *testing.T) {
	logger := zerolog.Nop()
	healthHandler := handlers.NewHealthHandler(nil)
	router := NewRouter(logger, "http://localhost:3000", healthHandler)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "healthy", body["status"])
}

func TestNewRouter_NotFound(t *testing.T) {
	logger := zerolog.Nop()
	healthHandler := handlers.NewHealthHandler(nil)
	router := NewRouter(logger, "http://localhost:3000", healthHandler)

	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestNewRouter_SetsRequestID(t *testing.T) {
	logger := zerolog.Nop()
	healthHandler := handlers.NewHealthHandler(nil)
	router := NewRouter(logger, "http://localhost:3000", healthHandler)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.NotEmpty(t, w.Header().Get("X-Request-ID"))
}
