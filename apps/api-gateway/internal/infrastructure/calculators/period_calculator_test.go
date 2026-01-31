package calculators

import (
	"testing"

	"github.com/melegattip/financial-resume-engine/internal/core/usecases"
	"github.com/stretchr/testify/assert"
)

func TestPeriodCalculator_FormatPeriodLabel(t *testing.T) {
	calc := NewPeriodCalculator()

	tests := []struct {
		name     string
		period   usecases.DatePeriod
		expected string
	}{
		{
			name: "year only",
			period: usecases.DatePeriod{
				Year: func() *int { y := 2024; return &y }(),
			},
			expected: "Año 2024",
		},
		{
			name: "year and month",
			period: usecases.DatePeriod{
				Year:  func() *int { y := 2024; return &y }(),
				Month: func() *int { m := 12; return &m }(),
			},
			expected: "Diciembre 2024",
		},
		{
			name: "month only",
			period: usecases.DatePeriod{
				Month: func() *int { m := 6; return &m }(),
			},
			expected: "Período inválido",
		},
		{
			name:     "empty period",
			period:   usecases.DatePeriod{},
			expected: "Todos los períodos",
		},
		{
			name: "january",
			period: usecases.DatePeriod{
				Year:  func() *int { y := 2024; return &y }(),
				Month: func() *int { m := 1; return &m }(),
			},
			expected: "Enero 2024",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.FormatPeriodLabel(tt.period)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewPeriodCalculator(t *testing.T) {
	calc := NewPeriodCalculator()
	assert.NotNil(t, calc)
}
