package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthResponse represents the JSON response for the health check endpoint.
type HealthResponse struct {
	Status    string                  `json:"status"`
	Timestamp string                  `json:"timestamp"`
	Checks    map[string]CheckResult  `json:"checks"`
	Version   string                  `json:"version"`
}

// CheckResult represents the health status of a single component.
type CheckResult struct {
	Status    string `json:"status"`
	LatencyMs *int64 `json:"latency_ms,omitempty"`
	Error     string `json:"error,omitempty"`
}

// HealthHandler handles the GET /health endpoint.
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler creates a new HealthHandler.
// db can be nil if database is not yet configured.
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Handle processes the health check request.
func (h *HealthHandler) Handle(c *gin.Context) {
	checks := make(map[string]CheckResult)
	overallStatus := "healthy"

	if h.db != nil {
		dbCheck := h.checkDatabase()
		checks["database"] = dbCheck
		if dbCheck.Status == "down" {
			overallStatus = "degraded"
		}
	}

	c.JSON(http.StatusOK, HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    checks,
		Version:   "0.1.0",
	})
}

func (h *HealthHandler) checkDatabase() CheckResult {
	sqlDB, err := h.db.DB()
	if err != nil {
		return CheckResult{Status: "down", Error: err.Error()}
	}

	start := time.Now()
	err = sqlDB.Ping()
	latency := time.Since(start).Milliseconds()

	if err != nil {
		return CheckResult{Status: "down", Error: err.Error()}
	}

	return CheckResult{Status: "up", LatencyMs: &latency}
}
