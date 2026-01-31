package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions/incomes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockIncomeService struct {
	mock.Mock
}

func (m *MockIncomeService) CreateIncome(ctx context.Context, request *incomes.CreateIncomeRequest) (*incomes.CreateIncomeResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*incomes.CreateIncomeResponse), args.Error(1)
}

func (m *MockIncomeService) GetIncome(ctx context.Context, userID string, incomeID string) (*incomes.GetIncomeResponse, error) {
	args := m.Called(ctx, userID, incomeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*incomes.GetIncomeResponse), args.Error(1)
}

func (m *MockIncomeService) ListIncomes(ctx context.Context, userID string) (*incomes.ListIncomesResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*incomes.ListIncomesResponse), args.Error(1)
}

func (m *MockIncomeService) UpdateIncome(ctx context.Context, userID string, incomeID string, request *incomes.UpdateIncomeRequest) (*incomes.UpdateIncomeResponse, error) {
	args := m.Called(ctx, userID, incomeID, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*incomes.UpdateIncomeResponse), args.Error(1)
}

func (m *MockIncomeService) DeleteIncome(ctx context.Context, userID string, incomeID string) error {
	args := m.Called(ctx, userID, incomeID)
	return args.Error(0)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestCreateIncome(t *testing.T) {
	mockService := new(MockIncomeService)
	handler := NewIncomeHandler(mockService)

	router := setupRouter()
	router.POST("/incomes", handler.CreateIncome)

	request := incomes.CreateIncomeRequest{
		UserID:      "user1",
		Amount:      1000.50,
		Description: "Salary",
		CategoryID:  "work",
	}

	expectedResponse := &incomes.CreateIncomeResponse{
		ID:          "1",
		UserID:      request.UserID,
		Amount:      request.Amount,
		Description: request.Description,
		CategoryID:  request.CategoryID,
		CreatedAt:   "2024-01-01T00:00:00Z",
		UpdatedAt:   "2024-01-01T00:00:00Z",
	}

	mockService.On("CreateIncome", mock.Anything, &request).Return(expectedResponse, nil)

	body, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/incomes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response incomes.CreateIncomeResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, expectedResponse.ID, response.ID)
	assert.Equal(t, expectedResponse.UserID, response.UserID)
	assert.Equal(t, expectedResponse.Amount, response.Amount)
	assert.Equal(t, expectedResponse.Description, response.Description)
	assert.Equal(t, expectedResponse.CategoryID, response.CategoryID)
	mockService.AssertExpectations(t)
}

func TestGetIncome(t *testing.T) {
	mockService := new(MockIncomeService)
	handler := NewIncomeHandler(mockService)

	router := setupRouter()
	router.GET("/incomes/:user_id/:id", handler.GetIncome)

	expectedResponse := &incomes.GetIncomeResponse{
		ID:          "1",
		UserID:      "user1",
		Amount:      1000.50,
		Description: "Salary",
		CategoryID:  "work",
		CreatedAt:   "2024-01-01T00:00:00Z",
		UpdatedAt:   "2024-01-01T00:00:00Z",
	}

	mockService.On("GetIncome", mock.Anything, "user1", "1").Return(expectedResponse, nil)

	req, _ := http.NewRequest("GET", "/incomes/user1/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response incomes.GetIncomeResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, expectedResponse.ID, response.ID)
	assert.Equal(t, expectedResponse.UserID, response.UserID)
	assert.Equal(t, expectedResponse.Amount, response.Amount)
	assert.Equal(t, expectedResponse.Description, response.Description)
	assert.Equal(t, expectedResponse.CategoryID, response.CategoryID)
	mockService.AssertExpectations(t)
}

func TestListIncomes(t *testing.T) {
	mockService := new(MockIncomeService)
	handler := NewIncomeHandler(mockService)

	router := setupRouter()
	router.GET("/incomes/:user_id", handler.ListIncomes)

	expectedResponse := &incomes.ListIncomesResponse{
		Incomes: []incomes.GetIncomeResponse{
			{
				ID:          "1",
				UserID:      "user1",
				Amount:      1000.50,
				Description: "Salary",
				CategoryID:  "work",
				CreatedAt:   "2024-01-01T00:00:00Z",
				UpdatedAt:   "2024-01-01T00:00:00Z",
			},
			{
				ID:          "2",
				UserID:      "user1",
				Amount:      2000.00,
				Description: "Bonus",
				CategoryID:  "work",
				CreatedAt:   "2024-01-01T00:00:00Z",
				UpdatedAt:   "2024-01-01T00:00:00Z",
			},
		},
	}

	mockService.On("ListIncomes", mock.Anything, "user1").Return(expectedResponse, nil)

	req, _ := http.NewRequest("GET", "/incomes/user1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response incomes.ListIncomesResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Len(t, response.Incomes, 2)
	assert.Equal(t, expectedResponse.Incomes[0].ID, response.Incomes[0].ID)
	assert.Equal(t, expectedResponse.Incomes[1].ID, response.Incomes[1].ID)
	mockService.AssertExpectations(t)
}

func TestUpdateIncome(t *testing.T) {
	mockService := new(MockIncomeService)
	handler := NewIncomeHandler(mockService)

	router := setupRouter()
	router.PUT("/incomes/:user_id/:id", handler.UpdateIncome)

	request := incomes.UpdateIncomeRequest{
		Amount:      1500.00,
		Description: "Updated Salary",
		CategoryID:  "work",
	}

	expectedResponse := &incomes.UpdateIncomeResponse{
		ID:          "1",
		UserID:      "user1",
		Amount:      request.Amount,
		Description: request.Description,
		CategoryID:  request.CategoryID,
		CreatedAt:   "2024-01-01T00:00:00Z",
		UpdatedAt:   "2024-01-01T00:00:00Z",
	}

	mockService.On("UpdateIncome", mock.Anything, "user1", "1", &request).Return(expectedResponse, nil)

	body, _ := json.Marshal(request)
	req, _ := http.NewRequest("PUT", "/incomes/user1/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response incomes.UpdateIncomeResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, expectedResponse.ID, response.ID)
	assert.Equal(t, expectedResponse.UserID, response.UserID)
	assert.Equal(t, expectedResponse.Amount, response.Amount)
	assert.Equal(t, expectedResponse.Description, response.Description)
	assert.Equal(t, expectedResponse.CategoryID, response.CategoryID)
	mockService.AssertExpectations(t)
}

func TestDeleteIncome(t *testing.T) {
	mockService := new(MockIncomeService)
	handler := NewIncomeHandler(mockService)

	router := setupRouter()
	router.DELETE("/incomes/:user_id/:id", handler.DeleteIncome)

	mockService.On("DeleteIncome", mock.Anything, "user1", "1").Return(nil)

	req, _ := http.NewRequest("DELETE", "/incomes/user1/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockService.AssertExpectations(t)
}

func TestCreateIncome_Validation(t *testing.T) {
	mockService := new(MockIncomeService)
	handler := NewIncomeHandler(mockService)

	router := setupRouter()
	router.POST("/incomes", handler.CreateIncome)

	tests := []struct {
		name           string
		request        incomes.CreateIncomeRequest
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Empty description",
			request: incomes.CreateIncomeRequest{
				UserID:      "user1",
				Amount:      1000.50,
				Description: "",
				CategoryID:  "work",
				Source:      "Company XYZ",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "description is required",
		},
		{
			name: "Zero amount",
			request: incomes.CreateIncomeRequest{
				UserID:      "user1",
				Amount:      0,
				Description: "Salary",
				CategoryID:  "work",
				Source:      "Company XYZ",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "amount must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest("POST", "/incomes", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]string
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Equal(t, tt.expectedError, response["error"])
		})
	}
}
