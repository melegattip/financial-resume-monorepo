package ports

import (
	"context"

	"github.com/gin-gonic/gin"
)

// HealthChecker provides health status for a component.
type HealthChecker interface {
	Check(ctx context.Context) HealthStatus
}

// HealthStatus represents the health of a single component.
type HealthStatus struct {
	Status    string `json:"status"`               // "up" or "down"
	LatencyMs *int64 `json:"latency_ms,omitempty"`
	Error     string `json:"error,omitempty"`
}

// Module defines the contract for pluggable business modules.
// Each module (auth, transactions, gamification, etc.) implements this.
type Module interface {
	// Name returns the module identifier.
	Name() string
	// RegisterRoutes adds the module's HTTP routes to the router.
	RegisterRoutes(router *gin.RouterGroup)
	// RegisterSubscribers registers event bus subscriptions.
	RegisterSubscribers(bus EventBus)
}
