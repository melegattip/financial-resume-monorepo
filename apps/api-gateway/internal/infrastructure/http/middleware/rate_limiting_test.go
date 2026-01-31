package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/services"
	"github.com/stretchr/testify/assert"
)

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		requests      int
		limit         int
		expectBlocked bool
		userID        string
	}{
		{
			name:          "dentro del límite",
			requests:      5,
			limit:         10,
			expectBlocked: false,
			userID:        "user1",
		},
		{
			name:          "exactamente en el límite",
			requests:      10,
			limit:         10,
			expectBlocked: false,
			userID:        "user2",
		},
		{
			name:          "excede el límite",
			requests:      12,
			limit:         10,
			expectBlocked: true,
			userID:        "user3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Crear servicio de rate limiting
			rateLimitService := services.NewInMemoryRateLimitService()

			// Configuración de rate limiting
			config := RateLimitConfig{
				RequestsPerMinute: tt.limit,
				SkipPaths:         []string{"/health"},
			}

			// Crear router con middleware
			router := gin.New()

			// Middleware mock para simular autenticación
			router.Use(func(c *gin.Context) {
				c.Set("user_id", tt.userID)
				c.Next()
			})

			router.Use(RateLimitMiddleware(rateLimitService, config))

			// Endpoint de prueba
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			var lastStatus int
			var blockedCount int

			// Simular múltiples requests
			for i := 0; i < tt.requests; i++ {
				req := httptest.NewRequest("GET", "/test", nil)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)
				lastStatus = w.Code

				if w.Code == http.StatusTooManyRequests {
					blockedCount++
				}

				// Verificar que requests dentro del límite sean permitidas
				if i < tt.limit {
					assert.Equal(t, http.StatusOK, w.Code, "Request %d should be allowed (within limit)", i+1)
				}
			}

			// Verificar el comportamiento general
			if tt.expectBlocked {
				assert.True(t, blockedCount > 0, "Should have blocked some requests")
				assert.Equal(t, http.StatusTooManyRequests, lastStatus, "Last request should be rate limited")
			} else {
				assert.Equal(t, 0, blockedCount, "Should not have blocked any requests")
				assert.Equal(t, http.StatusOK, lastStatus, "Last request should not be rate limited")
			}
		})
	}
}

func TestRateLimitMiddleware_SkipPaths(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rateLimitService := services.NewInMemoryRateLimitService()
	config := RateLimitConfig{
		RequestsPerMinute: 1, // Límite muy bajo para asegurar que paths normales serían bloqueados
		SkipPaths:         []string{"/health", "/swagger/"},
	}

	router := gin.New()

	// Middleware de autenticación mock para endpoints que lo necesitan
	router.Use(func(c *gin.Context) {
		path := c.Request.URL.Path
		// Solo establecer user_id para endpoints que no están en SkipPaths
		if !strings.HasPrefix(path, "/health") && !strings.HasPrefix(path, "/swagger/") {
			c.Set("user_id", "testuser")
		}
		c.Next()
	})

	router.Use(RateLimitMiddleware(rateLimitService, config))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/swagger/index.html", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"swagger": "docs"})
	})

	router.GET("/api/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	// Múltiples requests a paths que deben ser saltados
	for i := 0; i < 10; i++ {
		// Test /health - debe ser siempre permitido
		req := httptest.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "Health endpoint should never be rate limited")

		// Test /swagger/ - debe ser siempre permitido
		req = httptest.NewRequest("GET", "/swagger/index.html", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "Swagger endpoint should never be rate limited")
	}

	// Verificar que paths no incluidos en SkipPaths sí sean limitados
	// First request should pass
	req := httptest.NewRequest("GET", "/api/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "First request to /api/test should be allowed")

	// Second request should be blocked (limit = 1)
	req = httptest.NewRequest("GET", "/api/test", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code, "Second request to /api/test should be rate limited")
}

func TestRateLimitMiddleware_IPLimiting(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rateLimitService := services.NewInMemoryRateLimitService()
	config := RateLimitConfig{
		RequestsPerMinute:   100, // Alto para usuario
		IPLimitEnabled:      true,
		IPRequestsPerMinute: 3, // Bajo para IP
	}

	router := gin.New()

	// Middleware mock para simular autenticación
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "user1")
		c.Next()
	})

	router.Use(RateLimitMiddleware(rateLimitService, config))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	testIP := "192.168.1.100"

	// Hacer requests hasta exceder el límite de IP
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Forwarded-For", testIP)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Los primeros 3 requests deberían pasar (límite de IP = 3)
		if i < 3 {
			assert.Equal(t, http.StatusOK, w.Code, "Request %d should be allowed", i+1)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, w.Code, "Request %d should be rate limited by IP", i+1)
		}
	}
}

func TestRateLimitMiddleware_EndpointSpecificLimits(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rateLimitService := services.NewInMemoryRateLimitService()
	config := RateLimitConfig{
		RequestsPerMinute: 100, // Alto para uso general
		EndpointLimits: map[string]int{
			"POST:/api/v1/auth/login": 2, // Límite específico muy bajo para login
		},
	}

	router := gin.New()

	userID := "user1"

	// Middleware mock para simular autenticación
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})

	router.Use(RateLimitMiddleware(rateLimitService, config))

	router.POST("/api/v1/auth/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "login success"})
	})

	router.GET("/api/v1/dashboard", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "dashboard"})
	})

	// Test endpoint con límite específico (límite = 2)
	for i := 0; i < 4; i++ {
		req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(`{"email":"test@test.com"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if i < 2 {
			assert.Equal(t, http.StatusOK, w.Code, "Login request %d should be allowed", i+1)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, w.Code, "Login request %d should be rate limited", i+1)
		}
	}

	// Test endpoint sin límite específico (debería usar el límite general)
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/api/v1/dashboard", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "Dashboard request %d should be allowed", i+1)
	}
}

func TestAntiSpamMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rateLimitService := services.NewInMemoryRateLimitService()

	router := gin.New()

	userID := "spammer1"

	// Middleware mock para simular autenticación
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})

	router.Use(AntiSpamMiddleware(rateLimitService))

	router.POST("/api/v1/auth/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "login"})
	})

	router.GET("/api/v1/dashboard", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "dashboard"})
	})

	// Test: Exceder límite en endpoint sensible (límite = 5 para endpoints sensibles)
	for i := 0; i < 7; i++ {
		req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(`{"email":"test@test.com"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if i < 5 {
			assert.Equal(t, http.StatusOK, w.Code, "Login request %d should be allowed", i+1)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, w.Code, "Login request %d should be blocked by anti-spam", i+1)
		}
	}
}

func TestAntiSpamMiddleware_RapidRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rateLimitService := services.NewInMemoryRateLimitService()

	router := gin.New()

	userID := "rapiduser"

	// Middleware mock para simular autenticación
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})

	router.Use(AntiSpamMiddleware(rateLimitService))

	router.GET("/api/v1/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	// Simular requests rápidas (límite es 10 en 10 segundos para rate limiting rápido)
	// Los primeros 10 deberían pasar, luego el middleware detecta actividad sospechosa
	for i := 0; i < 12; i++ {
		req := httptest.NewRequest("GET", "/api/v1/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if i < 10 {
			assert.Equal(t, http.StatusOK, w.Code, "Request %d should be allowed", i+1)
		}
		// Después de 10 requests, el middleware registra actividad sospechosa
		// pero no necesariamente bloquea todas las requests inmediatamente
	}

	// Verificar que se registró actividad sospechosa después de exceder el límite
	suspiciousCount, err := rateLimitService.GetSuspiciousActivityCount(context.Background(), userID, "rapid_requests")
	assert.NoError(t, err)
	assert.True(t, suspiciousCount > 0, "Should have recorded suspicious activity")
}

func TestMetricsMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rateLimitService := services.NewInMemoryRateLimitService()

	router := gin.New()

	userID := "testuser"

	// Middleware mock para simular autenticación
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})

	router.Use(MetricsMiddleware(rateLimitService))

	router.GET("/api/v1/test", func(c *gin.Context) {
		time.Sleep(50 * time.Millisecond) // Simular procesamiento
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	router.GET("/api/v1/error", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "test error"})
	})

	// Test request exitosa
	req := httptest.NewRequest("GET", "/api/v1/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test request con error
	req = httptest.NewRequest("GET", "/api/v1/error", nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Esperar un poco para que las métricas se registren
	time.Sleep(100 * time.Millisecond)

	// Verificar que se registraron métricas
	// Nota: Debido a que Go maps no tienen orden garantizado, vamos a verificar
	// que el sistema de métricas esté funcionando de una manera más simple

	// Verificar que hay métricas de performance registradas (estas siempre se crean)
	performanceTags := map[string]string{
		"endpoint": "GET:/api/v1/test",
		"range":    "fast", // El test usa 50ms de sleep, debería ser "fast"
	}
	performanceMetrics, err := rateLimitService.GetMetric(context.Background(), "api_performance", "daily", performanceTags)
	assert.NoError(t, err)
	assert.True(t, performanceMetrics >= 1, "Should have recorded performance metrics")

	// Verificar métricas de usuario (estas son más simples)
	userTags := map[string]string{
		"user_id": userID,
	}
	userMetrics, err := rateLimitService.GetMetric(context.Background(), "user_requests", "daily", userTags)
	assert.NoError(t, err)
	assert.True(t, userMetrics >= 2, "Should have recorded user requests (2 calls made)")
}

func TestGetClientIP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		headers    map[string]string
		expectedIP string
	}{
		{
			name: "X-Forwarded-For header con múltiples IPs",
			headers: map[string]string{
				"X-Forwarded-For": "192.168.1.100, 10.0.0.1",
			},
			expectedIP: "192.168.1.100",
		},
		{
			name: "X-Forwarded-For header con una sola IP",
			headers: map[string]string{
				"X-Forwarded-For": "192.168.1.150",
			},
			expectedIP: "192.168.1.150",
		},
		{
			name: "X-Real-IP header",
			headers: map[string]string{
				"X-Real-IP": "192.168.1.200",
			},
			expectedIP: "192.168.1.200",
		},
		{
			name: "X-Forwarded-IP header",
			headers: map[string]string{
				"X-Forwarded-IP": "192.168.1.300",
			},
			expectedIP: "192.168.1.300",
		},
		{
			name: "Múltiples headers - X-Forwarded-For tiene prioridad",
			headers: map[string]string{
				"X-Forwarded-For": "192.168.1.100",
				"X-Real-IP":       "192.168.1.200",
				"X-Forwarded-IP":  "192.168.1.300",
			},
			expectedIP: "192.168.1.100",
		},
		{
			name: "Solo X-Real-IP cuando no hay X-Forwarded-For",
			headers: map[string]string{
				"X-Real-IP":      "192.168.1.200",
				"X-Forwarded-IP": "192.168.1.300",
			},
			expectedIP: "192.168.1.200",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			ip := getClientIP(c)
			assert.Equal(t, tt.expectedIP, ip, "Should extract correct client IP")
		})
	}

	// Test caso sin headers especiales
	t.Run("Fallback a ClientIP cuando no hay headers", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		ip := getClientIP(c)
		// El IP por defecto debería ser válido (aunque sea 127.0.0.1 en tests)
		assert.NotEmpty(t, ip, "Should return some IP address")
	})
}
