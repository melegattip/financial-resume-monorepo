package adapters

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/financial-ai-service/internal/adapters/cache"
	"github.com/financial-ai-service/internal/core/ports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// CacheClientTestSuite suite de tests para el cliente de cache
type CacheClientTestSuite struct {
	suite.Suite
	client ports.CacheClient
	ctx    context.Context
}

// SetupTest configura cada test individual
func (suite *CacheClientTestSuite) SetupTest() {
	// El cliente siempre usa mock por ahora
	suite.client = cache.NewRedisClient("redis://localhost:6379")
	suite.ctx = context.Background()
}

// TestNewRedisClient testa la creación del cliente Redis
func (suite *CacheClientTestSuite) TestNewRedisClient() {
	testCases := []struct {
		name string
		url  string
	}{
		{
			name: "URL estándar",
			url:  "redis://localhost:6379",
		},
		{
			name: "URL con autenticación",
			url:  "redis://user:pass@localhost:6379",
		},
		{
			name: "URL con base de datos",
			url:  "redis://localhost:6379/1",
		},
		{
			name: "URL vacía",
			url:  "",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Act
			client := cache.NewRedisClient(tc.url)

			// Assert
			assert.NotNil(suite.T(), client)
		})
	}
}

// TestGet_MockMode testa la obtención de valores en modo mock
func (suite *CacheClientTestSuite) TestGet_MockMode() {
	testCases := []struct {
		name        string
		key         string
		expectError bool
	}{
		{
			name:        "Clave simple",
			key:         "test:key",
			expectError: true, // Mock siempre retorna error "not found"
		},
		{
			name:        "Clave compleja",
			key:         "health_analysis:user123:monthly",
			expectError: true,
		},
		{
			name:        "Clave vacía",
			key:         "",
			expectError: true,
		},
		{
			name:        "Clave con caracteres especiales",
			key:         "user:123:data:$special",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Act
			result, err := suite.client.Get(suite.ctx, tc.key)

			// Assert
			if tc.expectError {
				assert.Error(suite.T(), err)
				assert.Nil(suite.T(), result)
				assert.Contains(suite.T(), err.Error(), "key not found in mock cache")
			} else {
				assert.NoError(suite.T(), err)
				assert.NotNil(suite.T(), result)
			}
		})
	}
}

// TestSet_MockMode testa el almacenamiento de valores en modo mock
func (suite *CacheClientTestSuite) TestSet_MockMode() {
	testCases := []struct {
		name  string
		key   string
		value []byte
		ttl   time.Duration
	}{
		{
			name:  "Datos simples",
			key:   "test:key",
			value: []byte(`{"test": "value"}`),
			ttl:   5 * time.Minute,
		},
		{
			name:  "Análisis de salud",
			key:   "health_analysis:user123",
			value: []byte(`{"score": 750, "level": "Bueno"}`),
			ttl:   20 * time.Hour,
		},
		{
			name:  "Insights",
			key:   "insights:user456",
			value: []byte(`[{"title": "Ahorro excelente"}]`),
			ttl:   24 * time.Hour,
		},
		{
			name:  "TTL corto",
			key:   "temp:data",
			value: []byte(`{"temp": true}`),
			ttl:   1 * time.Second,
		},
		{
			name:  "Sin TTL",
			key:   "permanent:data",
			value: []byte(`{"permanent": true}`),
			ttl:   0,
		},
		{
			name:  "Datos vacíos",
			key:   "empty:data",
			value: []byte{},
			ttl:   5 * time.Minute,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Act
			err := suite.client.Set(suite.ctx, tc.key, tc.value, tc.ttl)

			// Assert
			assert.NoError(suite.T(), err) // Mock siempre retorna éxito
		})
	}
}

// TestDelete_MockMode testa la eliminación de valores en modo mock
func (suite *CacheClientTestSuite) TestDelete_MockMode() {
	testCases := []struct {
		name string
		key  string
	}{
		{
			name: "Clave existente",
			key:  "existing:key",
		},
		{
			name: "Clave no existente",
			key:  "nonexistent:key",
		},
		{
			name: "Clave vacía",
			key:  "",
		},
		{
			name: "Clave con patrón",
			key:  "pattern:*:wildcard",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Act
			err := suite.client.Delete(suite.ctx, tc.key)

			// Assert
			assert.NoError(suite.T(), err) // Mock siempre retorna éxito
		})
	}
}

// TestClose_MockMode testa el cierre de conexión en modo mock
func (suite *CacheClientTestSuite) TestClose_MockMode() {
	// Act
	err := suite.client.Close()

	// Assert
	assert.NoError(suite.T(), err) // Mock siempre retorna éxito
}

// TestContextCancellation testa el manejo de cancelación de contexto
func (suite *CacheClientTestSuite) TestContextCancellation() {
	// Arrange
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel() // Cancelar inmediatamente

	// Act
	result, err := suite.client.Get(cancelledCtx, "test:key")

	// Assert
	// En modo mock, no debería importar la cancelación del contexto
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "key not found in mock cache")
}

// TestContextWithTimeout testa el manejo de timeout
func (suite *CacheClientTestSuite) TestContextWithTimeout() {
	// Arrange
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// Act
	err := suite.client.Set(ctxWithTimeout, "test:key", []byte("value"), 5*time.Minute)

	// Assert
	// En modo mock, debería funcionar independientemente del timeout
	assert.NoError(suite.T(), err)
}

// TestCacheOperations_SequentialFlow testa el flujo secuencial de operaciones
func (suite *CacheClientTestSuite) TestCacheOperations_SequentialFlow() {
	// Arrange
	key := "sequential:test"
	value := []byte(`{"sequential": true}`)
	ttl := 10 * time.Minute

	// Act & Assert - Set
	err := suite.client.Set(suite.ctx, key, value, ttl)
	assert.NoError(suite.T(), err)

	// Act & Assert - Get (debería fallar en mock)
	result, err := suite.client.Get(suite.ctx, key)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)

	// Act & Assert - Delete
	err = suite.client.Delete(suite.ctx, key)
	assert.NoError(suite.T(), err)

	// Act & Assert - Get después de delete (debería seguir fallando)
	result, err = suite.client.Get(suite.ctx, key)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

// TestConcurrentOperations testa operaciones concurrentes
func (suite *CacheClientTestSuite) TestConcurrentOperations() {
	// Arrange
	numOperations := 10
	errorsChan := make(chan error, numOperations*3) // Set, Get, Delete

	// Act
	for i := 0; i < numOperations; i++ {
		go func(index int) {
			key := fmt.Sprintf("concurrent:test:%d", index)
			value := []byte(fmt.Sprintf(`{"index": %d}`, index))

			// Set
			if err := suite.client.Set(suite.ctx, key, value, 5*time.Minute); err != nil {
				errorsChan <- err
			}

			// Get
			if _, err := suite.client.Get(suite.ctx, key); err == nil {
				// En mock, Get siempre falla, así que si no falla es un error
				errorsChan <- fmt.Errorf("expected error from mock Get operation")
			}

			// Delete
			if err := suite.client.Delete(suite.ctx, key); err != nil {
				errorsChan <- err
			}
		}(i)
	}

	// Assert - esperar un momento para que terminen las goroutines
	time.Sleep(100 * time.Millisecond)

	select {
	case err := <-errorsChan:
		suite.T().Errorf("Unexpected error in concurrent operations: %v", err)
	default:
		// No hay errores, está bien
	}
}

// TestDifferentDataTypes testa el almacenamiento de diferentes tipos de datos
func (suite *CacheClientTestSuite) TestDifferentDataTypes() {
	testCases := []struct {
		name string
		key  string
		data []byte
		ttl  time.Duration
	}{
		{
			name: "JSON objeto",
			key:  "json:object",
			data: []byte(`{"name": "test", "value": 123}`),
			ttl:  5 * time.Minute,
		},
		{
			name: "JSON array",
			key:  "json:array",
			data: []byte(`[1, 2, 3, "test"]`),
			ttl:  5 * time.Minute,
		},
		{
			name: "String simple",
			key:  "string:simple",
			data: []byte("simple string value"),
			ttl:  5 * time.Minute,
		},
		{
			name: "Datos binarios",
			key:  "binary:data",
			data: []byte{0x00, 0x01, 0x02, 0xFF},
			ttl:  5 * time.Minute,
		},
		{
			name: "String con caracteres especiales",
			key:  "string:special",
			data: []byte("Texto con acentos: ñáéíóú 🎉"),
			ttl:  5 * time.Minute,
		},
		{
			name: "JSON complejo",
			key:  "json:complex",
			data: []byte(`{
				"user": "test123",
				"analysis": {
					"score": 750,
					"insights": [
						{"title": "Ahorro", "value": 0.3}
					]
				},
				"timestamp": "2024-01-01T00:00:00Z"
			}`),
			ttl: 1 * time.Hour,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Act
			err := suite.client.Set(suite.ctx, tc.key, tc.data, tc.ttl)

			// Assert
			assert.NoError(suite.T(), err)
		})
	}
}

// TestDifferentTTLValues testa diferentes valores de TTL
func (suite *CacheClientTestSuite) TestDifferentTTLValues() {
	testCases := []struct {
		name string
		ttl  time.Duration
	}{
		{
			name: "TTL cero (sin expiración)",
			ttl:  0,
		},
		{
			name: "TTL muy corto",
			ttl:  1 * time.Millisecond,
		},
		{
			name: "TTL normal",
			ttl:  5 * time.Minute,
		},
		{
			name: "TTL largo",
			ttl:  24 * time.Hour,
		},
		{
			name: "TTL muy largo",
			ttl:  30 * 24 * time.Hour, // 30 días
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Arrange
			key := fmt.Sprintf("ttl:test:%s", tc.name)
			value := []byte(`{"ttl_test": true}`)

			// Act
			err := suite.client.Set(suite.ctx, key, value, tc.ttl)

			// Assert
			assert.NoError(suite.T(), err)
		})
	}
}

// TestLargeData testa el manejo de datos grandes
func (suite *CacheClientTestSuite) TestLargeData() {
	// Arrange
	largeData := make([]byte, 1024*1024) // 1MB
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	// Act
	err := suite.client.Set(suite.ctx, "large:data", largeData, 5*time.Minute)

	// Assert
	assert.NoError(suite.T(), err)

	// Try to get it back (should fail in mock)
	result, err := suite.client.Get(suite.ctx, "large:data")
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

// TestKeyPatterns testa diferentes patrones de claves
func (suite *CacheClientTestSuite) TestKeyPatterns() {
	testCases := []struct {
		name string
		key  string
	}{
		{
			name: "Patrón jerárquico",
			key:  "app:user:123:profile",
		},
		{
			name: "Patrón con timestamps",
			key:  "data:2024:01:01:analysis",
		},
		{
			name: "Patrón con UUID",
			key:  "session:550e8400-e29b-41d4-a716-446655440000",
		},
		{
			name: "Patrón con caracteres especiales",
			key:  "user:análisis:financiero",
		},
		{
			name: "Patrón muy largo",
			key:  strings.Repeat("very:long:key:pattern:", 10) + "end",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Act
			err := suite.client.Set(suite.ctx, tc.key, []byte(`{"pattern": "test"}`), 5*time.Minute)

			// Assert
			assert.NoError(suite.T(), err)
		})
	}
}

// TestClientReusability testa la reutilización del cliente
func (suite *CacheClientTestSuite) TestClientReusability() {
	// Perform multiple operations with the same client
	operations := 100

	for i := 0; i < operations; i++ {
		key := fmt.Sprintf("reuse:test:%d", i)
		value := []byte(fmt.Sprintf(`{"operation": %d}`, i))

		// Set
		err := suite.client.Set(suite.ctx, key, value, 1*time.Minute)
		assert.NoError(suite.T(), err)

		// Get
		_, err = suite.client.Get(suite.ctx, key)
		assert.Error(suite.T(), err) // Expected in mock mode

		// Delete
		err = suite.client.Delete(suite.ctx, key)
		assert.NoError(suite.T(), err)
	}
}

// TestRunSuite ejecuta todos los tests del suite
func TestCacheClientTestSuite(t *testing.T) {
	suite.Run(t, new(CacheClientTestSuite))
}

// Tests adicionales para casos específicos

// TestEmptyValues testa el manejo de valores vacíos
func TestEmptyValues(t *testing.T) {
	// Arrange
	client := cache.NewRedisClient("redis://localhost:6379")
	ctx := context.Background()

	testCases := []struct {
		name  string
		value []byte
	}{
		{
			name:  "Slice vacío",
			value: []byte{},
		},
		{
			name:  "Nil slice",
			value: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			err := client.Set(ctx, "empty:test", tc.value, 5*time.Minute)

			// Assert
			assert.NoError(t, err)
		})
	}
}

// TestMultipleClients testa múltiples instancias de clientes
func TestMultipleClients(t *testing.T) {
	// Arrange
	client1 := cache.NewRedisClient("redis://localhost:6379")
	client2 := cache.NewRedisClient("redis://localhost:6379")
	client3 := cache.NewRedisClient("redis://different:6379")
	ctx := context.Background()

	clients := []ports.CacheClient{client1, client2, client3}

	// Act & Assert
	for i, client := range clients {
		key := fmt.Sprintf("client:%d:test", i)
		value := []byte(fmt.Sprintf(`{"client": %d}`, i))

		err := client.Set(ctx, key, value, 5*time.Minute)
		assert.NoError(t, err)

		_, err = client.Get(ctx, key)
		assert.Error(t, err) // Expected in mock mode

		err = client.Close()
		assert.NoError(t, err)
	}
}
