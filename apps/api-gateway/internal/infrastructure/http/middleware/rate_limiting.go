package middleware

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	coreErrors "github.com/melegattip/financial-resume-engine/internal/core/errors"
	httpUtil "github.com/melegattip/financial-resume-engine/internal/infrastructure/http"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/services"
)

// RateLimitConfig configuración para rate limiting
type RateLimitConfig struct {
	RequestsPerMinute   int                                     // Requests por minuto general
	EndpointLimits      map[string]int                          // Límites específicos por endpoint
	IPLimitEnabled      bool                                    // Habilitar límite por IP
	IPRequestsPerMinute int                                     // Requests por minuto por IP
	SkipPaths           []string                                // Paths que no requieren rate limiting
	CustomLimitFunc     func(*gin.Context) (int, time.Duration) // Función personalizada para límites
}

// RateLimitMiddleware crea un middleware de rate limiting
func RateLimitMiddleware(rateLimitService *services.InMemoryRateLimitService, config RateLimitConfig) gin.HandlerFunc {
	// Configuración por defecto
	if config.RequestsPerMinute == 0 {
		config.RequestsPerMinute = 60 // 60 requests por minuto por defecto
	}
	if config.IPRequestsPerMinute == 0 {
		config.IPRequestsPerMinute = 120 // 120 requests por minuto por IP por defecto
	}
	if config.SkipPaths == nil {
		config.SkipPaths = []string{
			"/health",
			"/favicon.ico",
			"/robots.txt",
			"/manifest.json",
			"/swagger/",
			"/docs/",
		}
	}

	return gin.HandlerFunc(func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method

		// Verificar si el path debe ser saltado
		for _, skipPath := range config.SkipPaths {
			if strings.HasPrefix(path, skipPath) {
				c.Next()
				return
			}
		}

		// Obtener información del usuario y IP
		userID := getUserIDFromContext(c)
		clientIP := getClientIP(c)
		endpoint := method + ":" + path

		ctx := c.Request.Context()

		// 1. Verificar rate limit por IP (si está habilitado)
		if config.IPLimitEnabled && clientIP != "" {
			ipResult, err := rateLimitService.CheckIPRateLimit(ctx, clientIP, config.IPRequestsPerMinute, time.Minute)
			if err != nil {
				// Log error pero continuar con rate limit por usuario
				fmt.Printf("⚠️ Error verificando rate limit por IP: %v\n", err)
			} else if !ipResult.Allowed {
				// Rate limit por IP excedido
				setRateLimitHeaders(c, ipResult)
				httpUtil.TooManyRequests(c, coreErrors.NewTooManyRequests("Demasiadas requests desde esta IP. Intenta de nuevo más tarde."))
				c.Abort()
				return
			}
		}

		// 2. Verificar rate limit por usuario (si está autenticado)
		if userID != "" {
			// Determinar límite específico
			limit := config.RequestsPerMinute
			window := time.Minute

			// Verificar límite específico por endpoint
			if endpointLimit, exists := config.EndpointLimits[endpoint]; exists {
				limit = endpointLimit
			}

			// Usar función personalizada si está definida
			if config.CustomLimitFunc != nil {
				customLimit, customWindow := config.CustomLimitFunc(c)
				if customLimit > 0 {
					limit = customLimit
					window = customWindow
				}
			}

			// Verificar rate limit general por usuario
			userResult, err := rateLimitService.CheckRateLimit(ctx, userID, limit, window)
			if err != nil {
				fmt.Printf("⚠️ Error verificando rate limit por usuario: %v\n", err)
				// Continuar en caso de error
			} else if !userResult.Allowed {
				// Rate limit por usuario excedido
				setRateLimitHeaders(c, userResult)
				httpUtil.TooManyRequests(c, coreErrors.NewTooManyRequests("Demasiadas requests. Intenta de nuevo más tarde."))
				c.Abort()
				return
			}

			// Verificar rate limit específico por endpoint
			endpointResult, err := rateLimitService.CheckEndpointRateLimit(ctx, userID, endpoint, limit, window)
			if err != nil {
				fmt.Printf("⚠️ Error verificando rate limit por endpoint: %v\n", err)
			} else if !endpointResult.Allowed {
				setRateLimitHeaders(c, endpointResult)
				httpUtil.TooManyRequests(c, coreErrors.NewTooManyRequests(fmt.Sprintf("Demasiadas requests al endpoint %s. Intenta de nuevo más tarde.", endpoint)))
				c.Abort()
				return
			}

			// Agregar headers informativos si todo está bien
			setRateLimitHeaders(c, userResult)
		}

		// 3. Registrar métricas
		tags := map[string]string{
			"endpoint": endpoint,
			"user_id":  userID,
			"ip":       clientIP,
		}
		rateLimitService.IncrementMetric(ctx, "api_requests", tags)

		// Continuar con el siguiente middleware
		c.Next()
	})
}

// getUserIDFromContext extrae el user ID del contexto
func getUserIDFromContext(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		if userIDStr, ok := userID.(string); ok {
			return userIDStr
		}
	}
	return ""
}

// getClientIP obtiene la IP real del cliente considerando proxies
func getClientIP(c *gin.Context) string {
	// Verificar headers de proxy
	if ip := c.GetHeader("X-Forwarded-For"); ip != "" {
		// X-Forwarded-For puede contener múltiples IPs, tomar la primera
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0])
	}

	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		return ip
	}

	if ip := c.GetHeader("X-Forwarded-IP"); ip != "" {
		return ip
	}

	// Fallback a la IP de la conexión
	return c.ClientIP()
}

// setRateLimitHeaders establece headers estándar de rate limiting
func setRateLimitHeaders(c *gin.Context, result *services.RateLimitResult) {
	c.Header("X-RateLimit-Limit", strconv.FormatInt(result.Limit, 10))
	c.Header("X-RateLimit-Remaining", strconv.FormatInt(result.Remaining, 10))
	c.Header("X-RateLimit-Reset", strconv.FormatInt(result.ResetTime.Unix(), 10))

	if !result.Allowed {
		c.Header("Retry-After", strconv.FormatInt(int64(result.RetryAfter.Seconds()), 10))
	}
}

// AntiSpamMiddleware middleware especializado para detectar comportamiento sospechoso
func AntiSpamMiddleware(rateLimitService *services.InMemoryRateLimitService) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		userID := getUserIDFromContext(c)
		clientIP := getClientIP(c)
		path := c.Request.URL.Path

		if userID == "" {
			c.Next()
			return
		}

		ctx := c.Request.Context()

		// Detectar patrones sospechosos

		// 1. Verificar requests rápidas repetitivas (más de 10 requests en 10 segundos)
		quickLimit := 10
		quickWindow := 10 * time.Second
		quickResult, err := rateLimitService.CheckRateLimit(ctx, userID+":quick", quickLimit, quickWindow)
		if err == nil && !quickResult.Allowed {
			rateLimitService.RecordSuspiciousActivity(ctx, userID, "rapid_requests")
			fmt.Printf("🚨 Actividad sospechosa detectada - Requests rápidas: UserID=%s, IP=%s\n", userID, clientIP)
		}

		// 2. Detectar requests a endpoints sensibles (exceso de autenticación, etc.)
		sensitiveEndpoints := []string{
			"/api/v1/auth/login",
			"/api/v1/auth/register",
			"/api/v1/auth/refresh",
		}

		for _, sensitiveEndpoint := range sensitiveEndpoints {
			if strings.Contains(path, sensitiveEndpoint) {
				sensitiveLimit := 5
				sensitiveWindow := time.Minute
				sensitiveResult, err := rateLimitService.CheckEndpointRateLimit(ctx, userID, "sensitive:"+sensitiveEndpoint, sensitiveLimit, sensitiveWindow)
				if err == nil && !sensitiveResult.Allowed {
					rateLimitService.RecordSuspiciousActivity(ctx, userID, "auth_abuse")
					httpUtil.TooManyRequests(c, coreErrors.NewTooManyRequests("Demasiados intentos de autenticación. Cuenta temporalmente restringida."))
					c.Abort()
					return
				}
			}
		}

		// 3. Verificar acumulación de actividad sospechosa
		suspiciousCount, _ := rateLimitService.GetSuspiciousActivityCount(ctx, userID, "rapid_requests")
		authAbuseCount, _ := rateLimitService.GetSuspiciousActivityCount(ctx, userID, "auth_abuse")

		totalSuspicious := suspiciousCount + authAbuseCount
		if totalSuspicious >= 5 {
			// Usuario con demasiada actividad sospechosa - aplicar rate limiting más estricto
			strictLimit := 10
			strictWindow := 5 * time.Minute
			strictResult, err := rateLimitService.CheckRateLimit(ctx, userID+":restricted", strictLimit, strictWindow)
			if err == nil && !strictResult.Allowed {
				httpUtil.TooManyRequests(c, coreErrors.NewTooManyRequests("Cuenta temporalmente restringida por actividad sospechosa. Contacta soporte si crees que es un error."))
				c.Abort()
				return
			}
		}

		c.Next()
	})
}

// MetricsMiddleware middleware para registrar métricas de requests
func MetricsMiddleware(rateLimitService *services.InMemoryRateLimitService) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		start := time.Now()

		// Procesar request
		c.Next()

		// Registrar métricas después del procesamiento
		duration := time.Since(start)
		userID := getUserIDFromContext(c)
		endpoint := c.Request.Method + ":" + c.Request.URL.Path
		status := c.Writer.Status()

		ctx := c.Request.Context()

		// Métricas básicas
		rateLimitService.IncrementMetric(ctx, "api_requests_total", map[string]string{
			"endpoint": endpoint,
			"status":   strconv.Itoa(status),
		})

		// Métricas de performance (agrupar por rangos)
		var performanceRange string
		switch {
		case duration < 100*time.Millisecond:
			performanceRange = "fast"
		case duration < 500*time.Millisecond:
			performanceRange = "medium"
		case duration < 1*time.Second:
			performanceRange = "slow"
		default:
			performanceRange = "very_slow"
		}

		rateLimitService.IncrementMetric(ctx, "api_performance", map[string]string{
			"endpoint": endpoint,
			"range":    performanceRange,
		})

		// Métricas de errores
		if status >= 400 {
			errorType := "client_error"
			if status >= 500 {
				errorType = "server_error"
			}

			rateLimitService.IncrementMetric(ctx, "api_errors", map[string]string{
				"endpoint":   endpoint,
				"error_type": errorType,
				"status":     strconv.Itoa(status),
			})
		}

		// Métricas por usuario (solo si está autenticado)
		if userID != "" {
			rateLimitService.IncrementMetric(ctx, "user_requests", map[string]string{
				"user_id": userID,
			})
		}
	})
}
