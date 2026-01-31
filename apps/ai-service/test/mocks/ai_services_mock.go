package mocks

import (
	"context"
	"time"

	"github.com/financial-ai-service/internal/core/ports"
	"github.com/stretchr/testify/mock"
)

// MockAIAnalysisPort es un mock para el puerto de análisis de IA
type MockAIAnalysisPort struct {
	mock.Mock
}

func (m *MockAIAnalysisPort) AnalyzeFinancialHealth(ctx context.Context, data ports.FinancialAnalysisData) (*ports.HealthAnalysis, error) {
	args := m.Called(ctx, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ports.HealthAnalysis), args.Error(1)
}

func (m *MockAIAnalysisPort) GenerateInsights(ctx context.Context, data ports.FinancialAnalysisData) ([]ports.AIInsight, error) {
	args := m.Called(ctx, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ports.AIInsight), args.Error(1)
}

// MockPurchaseDecisionPort es un mock para el puerto de decisiones de compra
type MockPurchaseDecisionPort struct {
	mock.Mock
}

func (m *MockPurchaseDecisionPort) CanIBuy(ctx context.Context, request ports.PurchaseAnalysisRequest) (*ports.PurchaseDecision, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ports.PurchaseDecision), args.Error(1)
}

func (m *MockPurchaseDecisionPort) SuggestAlternatives(ctx context.Context, request ports.PurchaseAnalysisRequest) ([]ports.Alternative, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ports.Alternative), args.Error(1)
}

// MockCreditAnalysisPort es un mock para el puerto de análisis crediticio
type MockCreditAnalysisPort struct {
	mock.Mock
}

func (m *MockCreditAnalysisPort) GenerateImprovementPlan(ctx context.Context, data ports.FinancialAnalysisData) (*ports.CreditPlan, error) {
	args := m.Called(ctx, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ports.CreditPlan), args.Error(1)
}

func (m *MockCreditAnalysisPort) CalculateCreditScore(ctx context.Context, data ports.FinancialAnalysisData) (int, error) {
	args := m.Called(ctx, data)
	return args.Int(0), args.Error(1)
}

// MockOpenAIClient es un mock para el cliente de OpenAI
type MockOpenAIClient struct {
	mock.Mock
}

func (m *MockOpenAIClient) GenerateCompletion(ctx context.Context, prompt string) (string, error) {
	args := m.Called(ctx, prompt)
	return args.String(0), args.Error(1)
}

func (m *MockOpenAIClient) GenerateAnalysis(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	args := m.Called(ctx, systemPrompt, userPrompt)
	return args.String(0), args.Error(1)
}

// MockCacheClient es un mock para el cliente de cache
type MockCacheClient struct {
	mock.Mock
}

func (m *MockCacheClient) Get(ctx context.Context, key string) ([]byte, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockCacheClient) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *MockCacheClient) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCacheClient) Close() error {
	args := m.Called()
	return args.Error(0)
}
