package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/financial-ai-service/internal/core/ports"
)

// RedisClient implementa el adaptador de Redis para cache
type RedisClient struct {
	url      string
	password string
	db       int
	useMock  bool
}

// NewRedisClient crea un nuevo cliente de Redis
func NewRedisClient(url string) ports.CacheClient {
	// Por ahora usar mock hasta que se configure Redis
	log.Println("🎭 Redis client initialized in MOCK mode")

	return &RedisClient{
		url:     url,
		useMock: true,
	}
}

// Get obtiene un valor del cache
func (r *RedisClient) Get(ctx context.Context, key string) ([]byte, error) {
	if r.useMock {
		log.Printf("🎭 Mock cache GET: %s", key)
		return nil, fmt.Errorf("key not found in mock cache")
	}

	// TODO: Implementar Redis real
	return nil, fmt.Errorf("not implemented")
}

// Set almacena un valor en el cache
func (r *RedisClient) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if r.useMock {
		log.Printf("🎭 Mock cache SET: %s (TTL: %v)", key, ttl)
		return nil
	}

	// TODO: Implementar Redis real
	return fmt.Errorf("not implemented")
}

// Delete elimina un valor del cache
func (r *RedisClient) Delete(ctx context.Context, key string) error {
	if r.useMock {
		log.Printf("🎭 Mock cache DELETE: %s", key)
		return nil
	}

	// TODO: Implementar Redis real
	return fmt.Errorf("not implemented")
}

// Close cierra la conexión al cache
func (r *RedisClient) Close() error {
	if r.useMock {
		log.Println("🎭 Mock cache connection closed")
		return nil
	}

	// TODO: Implementar Redis real
	return fmt.Errorf("not implemented")
}
