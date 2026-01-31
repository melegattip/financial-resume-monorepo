package services

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"
)

// CacheEntry representa una entrada en el cache con TTL
type CacheEntry struct {
	Value     interface{}
	ExpiresAt time.Time
}

// InMemoryRateLimitService implementa rate limiting usando memoria local
type InMemoryRateLimitService struct {
	counters sync.Map // map[string]*CacheEntry
	metrics  sync.Map // map[string]*CacheEntry
}

// NewInMemoryRateLimitService crea una nueva instancia del servicio
func NewInMemoryRateLimitService() *InMemoryRateLimitService {
	service := &InMemoryRateLimitService{}

	// Iniciar goroutine para limpiar entradas expiradas
	go service.cleanupExpiredEntries()

	return service
}

// RateLimitResult representa el resultado de una verificación de rate limit
type RateLimitResult struct {
	Allowed    bool          `json:"allowed"`
	Count      int64         `json:"count"`
	Limit      int64         `json:"limit"`
	Remaining  int64         `json:"remaining"`
	ResetTime  time.Time     `json:"reset_time"`
	RetryAfter time.Duration `json:"retry_after"`
}

// CheckRateLimit verifica y actualiza el contador de rate limit para un usuario
func (s *InMemoryRateLimitService) CheckRateLimit(ctx context.Context, userID string, limit int, window time.Duration) (*RateLimitResult, error) {
	key := fmt.Sprintf("rate_limit:%s", userID)
	return s.checkRateLimitGeneric(key, limit, window)
}

// CheckIPRateLimit verifica rate limit por IP (protección adicional)
func (s *InMemoryRateLimitService) CheckIPRateLimit(ctx context.Context, ip string, limit int, window time.Duration) (*RateLimitResult, error) {
	key := fmt.Sprintf("rate_limit_ip:%s", ip)
	return s.checkRateLimitGeneric(key, limit, window)
}

// CheckEndpointRateLimit verifica rate limit por endpoint específico
func (s *InMemoryRateLimitService) CheckEndpointRateLimit(ctx context.Context, userID, endpoint string, limit int, window time.Duration) (*RateLimitResult, error) {
	key := fmt.Sprintf("rate_limit_endpoint:%s:%s", userID, endpoint)
	return s.checkRateLimitGeneric(key, limit, window)
}

// checkRateLimitGeneric función genérica para rate limiting
func (s *InMemoryRateLimitService) checkRateLimitGeneric(key string, limit int, window time.Duration) (*RateLimitResult, error) {
	now := time.Now()
	expiresAt := now.Add(window)

	// Obtener o crear contador
	var count int64 = 1
	if val, exists := s.counters.Load(key); exists {
		if entry, ok := val.(*CacheEntry); ok && entry.ExpiresAt.After(now) {
			// Incrementar contador existente
			if currentCount, ok := entry.Value.(int64); ok {
				count = currentCount + 1
				entry.Value = count
			}
		} else {
			// Contador expirado o inválido, crear nuevo
			s.counters.Store(key, &CacheEntry{
				Value:     int64(1),
				ExpiresAt: expiresAt,
			})
			count = 1
		}
	} else {
		// Crear nuevo contador
		s.counters.Store(key, &CacheEntry{
			Value:     int64(1),
			ExpiresAt: expiresAt,
		})
		count = 1
	}

	// Actualizar el contador
	s.counters.Store(key, &CacheEntry{
		Value:     count,
		ExpiresAt: expiresAt,
	})

	remaining := int64(limit) - count
	if remaining < 0 {
		remaining = 0
	}

	return &RateLimitResult{
		Allowed:    count <= int64(limit),
		Count:      count,
		Limit:      int64(limit),
		Remaining:  remaining,
		ResetTime:  expiresAt,
		RetryAfter: time.Until(expiresAt),
	}, nil
}

// RecordSuspiciousActivity registra actividad sospechosa
func (s *InMemoryRateLimitService) RecordSuspiciousActivity(ctx context.Context, userID, activityType string) error {
	key := fmt.Sprintf("suspicious:%s:%s", userID, activityType)
	expiresAt := time.Now().Add(time.Hour)

	var count int64 = 1
	if val, exists := s.counters.Load(key); exists {
		if entry, ok := val.(*CacheEntry); ok && entry.ExpiresAt.After(time.Now()) {
			if currentCount, ok := entry.Value.(int64); ok {
				count = currentCount + 1
			}
		}
	}

	s.counters.Store(key, &CacheEntry{
		Value:     count,
		ExpiresAt: expiresAt,
	})

	return nil
}

// GetSuspiciousActivityCount obtiene el contador de actividad sospechosa
func (s *InMemoryRateLimitService) GetSuspiciousActivityCount(ctx context.Context, userID, activityType string) (int64, error) {
	key := fmt.Sprintf("suspicious:%s:%s", userID, activityType)

	if val, exists := s.counters.Load(key); exists {
		if entry, ok := val.(*CacheEntry); ok && entry.ExpiresAt.After(time.Now()) {
			if count, ok := entry.Value.(int64); ok {
				return count, nil
			}
		}
	}

	return 0, nil
}

// IncrementMetric incrementa una métrica específica
func (s *InMemoryRateLimitService) IncrementMetric(ctx context.Context, metricName string, tags map[string]string) error {
	// Crear clave estable con tags ordenados para evitar no determinismo
	key := buildMetricKey(metricName, tags)

	// Incrementar métrica diaria
	dailyKey := fmt.Sprintf("%s:daily:%s", key, time.Now().Format("2006-01-02"))
	hourlyKey := fmt.Sprintf("%s:hourly:%s", key, time.Now().Format("2006-01-02:15"))

	s.incrementMetricKey(dailyKey, 7*24*time.Hour)
	s.incrementMetricKey(hourlyKey, 48*time.Hour)

	return nil
}

// incrementMetricKey incrementa una métrica específica con TTL
func (s *InMemoryRateLimitService) incrementMetricKey(key string, ttl time.Duration) {
	expiresAt := time.Now().Add(ttl)
	var count int64 = 1

	if val, exists := s.metrics.Load(key); exists {
		if entry, ok := val.(*CacheEntry); ok && entry.ExpiresAt.After(time.Now()) {
			if currentCount, ok := entry.Value.(int64); ok {
				count = currentCount + 1
			}
		}
	}

	s.metrics.Store(key, &CacheEntry{
		Value:     count,
		ExpiresAt: expiresAt,
	})
}

// GetMetric obtiene el valor de una métrica
func (s *InMemoryRateLimitService) GetMetric(ctx context.Context, metricName string, period string, tags map[string]string) (int64, error) {
	key := buildMetricKey(metricName, tags)

	var finalKey string
	switch period {
	case "daily":
		finalKey = fmt.Sprintf("%s:daily:%s", key, time.Now().Format("2006-01-02"))
	case "hourly":
		finalKey = fmt.Sprintf("%s:hourly:%s", key, time.Now().Format("2006-01-02:15"))
	default:
		finalKey = key
	}

	if val, exists := s.metrics.Load(finalKey); exists {
		if entry, ok := val.(*CacheEntry); ok && entry.ExpiresAt.After(time.Now()) {
			if count, ok := entry.Value.(int64); ok {
				return count, nil
			}
		}
	}

	return 0, nil
}

// buildMetricKey construye una clave determinística para métricas a partir del nombre y tags
func buildMetricKey(metricName string, tags map[string]string) string {
	key := fmt.Sprintf("metric:%s", metricName)
	if len(tags) == 0 {
		return key
	}
	// Ordenar las claves de tags para asegurar orden estable
	keys := make([]string, 0, len(tags))
	for k := range tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		key += fmt.Sprintf(":%s:%s", k, tags[k])
	}
	return key
}

// cleanupExpiredEntries limpia las entradas expiradas periódicamente
func (s *InMemoryRateLimitService) cleanupExpiredEntries() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()

		// Limpiar contadores expirados
		s.counters.Range(func(key, value interface{}) bool {
			if entry, ok := value.(*CacheEntry); ok {
				if entry.ExpiresAt.Before(now) {
					s.counters.Delete(key)
				}
			}
			return true
		})

		// Limpiar métricas expiradas
		s.metrics.Range(func(key, value interface{}) bool {
			if entry, ok := value.(*CacheEntry); ok {
				if entry.ExpiresAt.Before(now) {
					s.metrics.Delete(key)
				}
			}
			return true
		})
	}
}

// Ping verifica la conectividad (siempre ok para in-memory)
func (s *InMemoryRateLimitService) Ping(ctx context.Context) error {
	return nil
}

// Close no hace nada para in-memory
func (s *InMemoryRateLimitService) Close() error {
	return nil
}
