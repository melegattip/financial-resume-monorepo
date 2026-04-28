package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/ai/domain"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/ai/service"
)

func init() { gin.SetMode(gin.TestMode) }

func newTestAIHandler() *AIHandler {
	return &AIHandler{
		monthlyCache:   make(map[string]domain.MonthlyCoachingReport),
		emailSentMap:   make(map[string]bool),
		educationCache: make(map[string]domain.EducationContent),
		insightsCache:  make(map[string]cachedInsights),
	}
}

// ---------------------------------------------------------------------------
// educationCacheKey
// ---------------------------------------------------------------------------

func TestEducationCacheKey_SameWeek(t *testing.T) {
	t1 := time.Date(2026, 3, 16, 0, 0, 0, 0, time.UTC) // Monday
	t2 := time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC) // Wednesday
	assert.Equal(t, educationCacheKey("user1", t1), educationCacheKey("user1", t2))
}

func TestEducationCacheKey_DifferentWeek(t *testing.T) {
	t1 := time.Date(2026, 3, 16, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2026, 3, 23, 0, 0, 0, 0, time.UTC)
	assert.NotEqual(t, educationCacheKey("user1", t1), educationCacheKey("user1", t2))
}

func TestEducationCacheKey_DifferentUsers(t *testing.T) {
	t1 := time.Date(2026, 3, 16, 0, 0, 0, 0, time.UTC)
	assert.NotEqual(t, educationCacheKey("userA", t1), educationCacheKey("userB", t1))
}

// ---------------------------------------------------------------------------
// HandleMonthlyCoaching — validation
// ---------------------------------------------------------------------------

func TestHandleMonthlyCoaching_InvalidMonthFormat(t *testing.T) {
	h := newTestAIHandler()
	router := gin.New()
	router.POST("/ai/monthly-coaching", h.HandleMonthlyCoaching)

	body, _ := json.Marshal(domain.MonthlyCoachingRequest{PreviousMonth: "not-a-date"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/ai/monthly-coaching", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Contains(t, resp["error"], "YYYY-MM")
}

func TestHandleMonthlyCoaching_InvalidMonthMonth(t *testing.T) {
	h := newTestAIHandler()
	router := gin.New()
	router.POST("/ai/monthly-coaching", h.HandleMonthlyCoaching)

	body, _ := json.Marshal(domain.MonthlyCoachingRequest{PreviousMonth: "2026-13"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/ai/monthly-coaching", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleMonthlyCoaching_InvalidJSON(t *testing.T) {
	h := newTestAIHandler()
	router := gin.New()
	router.POST("/ai/monthly-coaching", h.HandleMonthlyCoaching)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/ai/monthly-coaching", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ---------------------------------------------------------------------------
// HandleMonthlyCoaching — cache hit
// ---------------------------------------------------------------------------

func TestHandleMonthlyCoaching_CacheHit(t *testing.T) {
	h := newTestAIHandler()
	cachedReport := domain.MonthlyCoachingReport{
		Month:     "2026-02",
		Sentiment: "positivo",
	}
	h.monthlyCache["user1_2026-02"] = cachedReport

	router := gin.New()
	router.POST("/ai/monthly-coaching", func(c *gin.Context) {
		c.Set("user_id", "user1")
		h.HandleMonthlyCoaching(c)
	})

	body, _ := json.Marshal(domain.MonthlyCoachingRequest{
		FinancialData: domain.FinancialAnalysisData{UserID: "user1"},
		PreviousMonth: "2026-02",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/ai/monthly-coaching", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, true, resp["cached"])
	reportMap := resp["report"].(map[string]interface{})
	assert.Equal(t, "positivo", reportMap["sentiment"])
}

// ---------------------------------------------------------------------------
// HandleMonthlyCoaching — force bypass cache
// ---------------------------------------------------------------------------

func TestHandleMonthlyCoaching_ForceBypassesCache(t *testing.T) {
	svc := service.NewAnalysisService(service.NewOpenAIClient("")) // mock mode
	h := &AIHandler{
		analysisService: svc,
		monthlyCache:    make(map[string]domain.MonthlyCoachingReport),
		emailSentMap:    make(map[string]bool),
		educationCache:  make(map[string]domain.EducationContent),
		insightsCache:   make(map[string]cachedInsights),
	}

	// Pre-populate the cache with a stale "positivo" report.
	h.monthlyCache["user1_2026-02"] = domain.MonthlyCoachingReport{
		Month:     "2026-02",
		Sentiment: "positivo",
	}

	router := gin.New()
	router.POST("/ai/monthly-coaching", func(c *gin.Context) {
		c.Set("user_id", "user1")
		h.HandleMonthlyCoaching(c)
	})

	body, _ := json.Marshal(domain.MonthlyCoachingRequest{
		FinancialData: domain.FinancialAnalysisData{UserID: "user1"},
		PreviousMonth: "2026-02",
		Force:         true,
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/ai/monthly-coaching", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	// Force=true must bypass the cache: cached field must be false.
	assert.Equal(t, false, resp["cached"])
}

// ---------------------------------------------------------------------------
// HandleEducationCards — validation and cache
// ---------------------------------------------------------------------------

func TestHandleEducationCards_InvalidJSON(t *testing.T) {
	h := newTestAIHandler()
	router := gin.New()
	router.POST("/ai/education-cards", h.HandleEducationCards)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/ai/education-cards", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleEducationCards_CacheHit(t *testing.T) {
	h := newTestAIHandler()
	now := time.Now().UTC()
	cacheKey := educationCacheKey("user1", now)
	h.educationCache[cacheKey] = domain.EducationContent{
		Cards:       []domain.EducationCard{{Topic: "ahorro", Title: "Fondo"}},
		GeneratedAt: now,
	}

	router := gin.New()
	router.POST("/ai/education-cards", func(c *gin.Context) {
		c.Set("user_id", "user1")
		h.HandleEducationCards(c)
	})

	body, _ := json.Marshal(domain.EducationRequest{
		FinancialData: domain.FinancialAnalysisData{UserID: "user1"},
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/ai/education-cards", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, true, resp["cached"])
	cards := resp["cards"].([]interface{})
	assert.Len(t, cards, 1)
}
