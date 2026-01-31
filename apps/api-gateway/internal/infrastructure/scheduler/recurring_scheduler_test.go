package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/ports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRecurringTransactionUseCase is a mock implementation
type MockRecurringTransactionUseCase struct {
	mock.Mock
}

func (m *MockRecurringTransactionUseCase) ProcessPendingTransactions(ctx context.Context) (*ports.BatchProcessResult, error) {
	args := m.Called(ctx)
	return args.Get(0).(*ports.BatchProcessResult), args.Error(1)
}

func (m *MockRecurringTransactionUseCase) SendPendingNotifications(ctx context.Context) (*ports.NotificationResult, error) {
	args := m.Called(ctx)
	return args.Get(0).(*ports.NotificationResult), args.Error(1)
}

// Implement other required methods as no-ops for testing
func (m *MockRecurringTransactionUseCase) CreateRecurringTransaction(ctx context.Context, request *ports.CreateRecurringTransactionRequest) (*ports.RecurringTransactionResponse, error) {
	return nil, nil
}

func (m *MockRecurringTransactionUseCase) GetRecurringTransaction(ctx context.Context, userID, transactionID string) (*ports.RecurringTransactionResponse, error) {
	return nil, nil
}

func (m *MockRecurringTransactionUseCase) ListRecurringTransactions(ctx context.Context, userID string, filters ports.RecurringTransactionFilters) (*ports.ListRecurringTransactionsResponse, error) {
	return nil, nil
}

func (m *MockRecurringTransactionUseCase) UpdateRecurringTransaction(ctx context.Context, userID, transactionID string, request *ports.UpdateRecurringTransactionRequest) (*ports.RecurringTransactionResponse, error) {
	return nil, nil
}

func (m *MockRecurringTransactionUseCase) DeleteRecurringTransaction(ctx context.Context, userID, transactionID string) error {
	return nil
}

func (m *MockRecurringTransactionUseCase) PauseRecurringTransaction(ctx context.Context, userID, transactionID string) error {
	return nil
}

func (m *MockRecurringTransactionUseCase) ResumeRecurringTransaction(ctx context.Context, userID, transactionID string) error {
	return nil
}

func (m *MockRecurringTransactionUseCase) ExecuteRecurringTransaction(ctx context.Context, userID, transactionID string) (*ports.ExecutionResult, error) {
	return nil, nil
}

func (m *MockRecurringTransactionUseCase) GetRecurringTransactionsDashboard(ctx context.Context, userID string) (*ports.RecurringDashboardResponse, error) {
	return nil, nil
}

func (m *MockRecurringTransactionUseCase) GetCashFlowProjection(ctx context.Context, userID string, months int) (*ports.CashFlowProjectionResponse, error) {
	return nil, nil
}

func TestNewRecurringTransactionScheduler(t *testing.T) {
	mockUseCase := &MockRecurringTransactionUseCase{}
	interval := 1 * time.Hour

	scheduler := NewRecurringTransactionScheduler(mockUseCase, interval)

	assert.NotNil(t, scheduler)
	assert.Equal(t, interval, scheduler.interval)
	assert.False(t, scheduler.IsRunning())
}

func TestSchedulerStartStop(t *testing.T) {
	mockUseCase := &MockRecurringTransactionUseCase{}

	// Mock the expected calls
	mockUseCase.On("ProcessPendingTransactions", mock.Anything).Return(&ports.BatchProcessResult{
		ProcessedCount: 0,
		SuccessCount:   0,
		FailureCount:   0,
	}, nil)

	mockUseCase.On("SendPendingNotifications", mock.Anything).Return(&ports.NotificationResult{
		SentCount:    0,
		FailureCount: 0,
	}, nil)

	scheduler := NewRecurringTransactionScheduler(mockUseCase, 100*time.Millisecond)

	// Test start
	scheduler.Start()
	assert.True(t, scheduler.IsRunning())

	// Wait a bit to let it process
	time.Sleep(150 * time.Millisecond)

	// Test stop
	scheduler.Stop()
	assert.False(t, scheduler.IsRunning())

	// Verify mocks were called
	mockUseCase.AssertExpectations(t)
}

func TestSchedulerGetStatus(t *testing.T) {
	mockUseCase := &MockRecurringTransactionUseCase{}
	interval := 1 * time.Hour

	scheduler := NewRecurringTransactionScheduler(mockUseCase, interval)

	status := scheduler.GetStatus()

	assert.Contains(t, status, "running")
	assert.Contains(t, status, "interval")
	assert.Equal(t, false, status["running"])
	assert.Equal(t, interval.String(), status["interval"])
}

func TestSchedulerProcessNow(t *testing.T) {
	mockUseCase := &MockRecurringTransactionUseCase{}

	// Mock the expected calls
	mockUseCase.On("ProcessPendingTransactions", mock.Anything).Return(&ports.BatchProcessResult{
		ProcessedCount: 1,
		SuccessCount:   1,
		FailureCount:   0,
	}, nil)

	mockUseCase.On("SendPendingNotifications", mock.Anything).Return(&ports.NotificationResult{
		SentCount:    0,
		FailureCount: 0,
	}, nil)

	scheduler := NewRecurringTransactionScheduler(mockUseCase, 1*time.Hour)

	// Test manual processing
	scheduler.ProcessNow()

	// Verify mocks were called
	mockUseCase.AssertExpectations(t)
}
