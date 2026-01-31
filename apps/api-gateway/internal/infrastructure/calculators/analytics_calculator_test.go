package calculators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnalyticsCalculator_CalculatePercentages(t *testing.T) {
	calc := NewAnalyticsCalculator()

	tests := []struct {
		name     string
		amount   float64
		total    float64
		expected float64
	}{
		{
			name:     "normal percentage",
			amount:   50.0,
			total:    100.0,
			expected: 50.0,
		},
		{
			name:     "zero total",
			amount:   50.0,
			total:    0.0,
			expected: 0.0,
		},
		{
			name:     "zero amount",
			amount:   0.0,
			total:    100.0,
			expected: 0.0,
		},
		{
			name:     "decimal result",
			amount:   33.33,
			total:    100.0,
			expected: 33.33,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.CalculatePercentages(tt.amount, tt.total)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAnalyticsCalculator_CalculateAverage(t *testing.T) {
	calc := NewAnalyticsCalculator()

	tests := []struct {
		name     string
		total    float64
		count    int
		expected float64
	}{
		{
			name:     "normal average",
			total:    100.0,
			count:    4,
			expected: 25.0,
		},
		{
			name:     "zero count",
			total:    100.0,
			count:    0,
			expected: 0.0,
		},
		{
			name:     "zero total",
			total:    0.0,
			count:    5,
			expected: 0.0,
		},
		{
			name:     "decimal result",
			total:    100.0,
			count:    3,
			expected: 33.333333333333336,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.CalculateAverage(tt.total, tt.count)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAnalyticsCalculator_GenerateColorSeed(t *testing.T) {
	calc := NewAnalyticsCalculator()

	tests := []struct {
		name       string
		identifier string
		minValue   int
		maxValue   int
	}{
		{
			name:       "food category",
			identifier: "food",
			minValue:   0,
			maxValue:   4294967295, // uint32 max
		},
		{
			name:       "transport category",
			identifier: "transport",
			minValue:   0,
			maxValue:   4294967295,
		},
		{
			name:       "empty string",
			identifier: "",
			minValue:   0,
			maxValue:   4294967295,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.GenerateColorSeed(tt.identifier)
			assert.GreaterOrEqual(t, result, tt.minValue)
			assert.LessOrEqual(t, result, tt.maxValue)
		})
	}
}

func TestAnalyticsCalculator_GenerateColorSeed_Consistency(t *testing.T) {
	calc := NewAnalyticsCalculator()

	// El mismo identificador debería generar el mismo seed
	identifier := "test-category"
	seed1 := calc.GenerateColorSeed(identifier)
	seed2 := calc.GenerateColorSeed(identifier)

	assert.Equal(t, seed1, seed2, "El mismo identificador debería generar el mismo color seed")
}

func TestNewAnalyticsCalculator(t *testing.T) {
	calc := NewAnalyticsCalculator()
	assert.NotNil(t, calc)
}
