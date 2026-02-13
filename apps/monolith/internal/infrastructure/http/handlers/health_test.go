package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupRouter(handler *HealthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/health", handler.Handle)
	return r
}

func TestHealthHandler_WithoutDB(t *testing.T) {
	handler := NewHealthHandler(nil)
	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, "healthy", resp.Status)
	assert.NotEmpty(t, resp.Timestamp)
	assert.Equal(t, "0.1.0", resp.Version)
	assert.Empty(t, resp.Checks)
}

func TestHealthHandler_ResponseFormat(t *testing.T) {
	handler := NewHealthHandler(nil)
	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var raw map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &raw)
	require.NoError(t, err)

	assert.Contains(t, raw, "status")
	assert.Contains(t, raw, "timestamp")
	assert.Contains(t, raw, "checks")
	assert.Contains(t, raw, "version")
}

func openGormWithMock(t *testing.T, opts ...func(*sqlmock.Sqlmock)) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	// Without MonitorPingsOption, pings auto-succeed (good for init)
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	// GORM postgres driver runs SELECT VERSION() during Open
	mock.ExpectQuery("SELECT VERSION()").
		WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("PostgreSQL 15.0"))

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: mockDB,
	}), &gorm.Config{})
	require.NoError(t, err)

	t.Cleanup(func() { mockDB.Close() })
	return gormDB, mock
}

func TestHealthHandler_WithDB_Healthy(t *testing.T) {
	// Without MonitorPingsOption, pings always succeed
	gormDB, _ := openGormWithMock(t)

	handler := NewHealthHandler(gormDB)
	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, "healthy", resp.Status)
	assert.Contains(t, resp.Checks, "database")
	assert.Equal(t, "up", resp.Checks["database"].Status)
	assert.NotNil(t, resp.Checks["database"].LatencyMs)
}

func TestHealthHandler_WithDB_Down(t *testing.T) {
	// Use MonitorPingsOption so we can control ping failure
	mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)
	defer mockDB.Close()

	// Init expectations: ping + version query
	mock.ExpectPing()
	mock.ExpectQuery("SELECT VERSION()").
		WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("PostgreSQL 15.0"))

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: mockDB,
	}), &gorm.Config{})
	require.NoError(t, err)

	// Health check ping should fail
	mock.ExpectPing().WillReturnError(assert.AnError)

	handler := NewHealthHandler(gormDB)
	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp HealthResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, "degraded", resp.Status)
	assert.Contains(t, resp.Checks, "database")
	assert.Equal(t, "down", resp.Checks["database"].Status)
	assert.NotEmpty(t, resp.Checks["database"].Error)
}
