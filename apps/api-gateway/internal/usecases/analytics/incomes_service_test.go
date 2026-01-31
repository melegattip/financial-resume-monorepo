package analytics

import (
	"testing"

	"github.com/melegattip/financial-resume-engine/internal/core/usecases"
	"github.com/stretchr/testify/assert"
)

func TestIncomesAnalyticsService_validateParams(t *testing.T) {
	service := &IncomesAnalyticsService{}

	tests := []struct {
		name    string
		params  usecases.IncomesSummaryParams
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid params",
			params: usecases.IncomesSummaryParams{
				UserID: "user-123",
				Period: usecases.DatePeriod{},
			},
			wantErr: false,
		},
		{
			name: "empty user ID",
			params: usecases.IncomesSummaryParams{
				UserID: "",
			},
			wantErr: true,
			errMsg:  "El ID del usuario es requerido",
		},
		{
			name: "invalid year",
			params: usecases.IncomesSummaryParams{
				UserID: "user-123",
				Period: usecases.DatePeriod{
					Year: func() *int { y := 1800; return &y }(),
				},
			},
			wantErr: true,
			errMsg:  "Año inválido",
		},
		{
			name: "invalid month",
			params: usecases.IncomesSummaryParams{
				UserID: "user-123",
				Period: usecases.DatePeriod{
					Month: func() *int { m := 13; return &m }(),
				},
			},
			wantErr: true,
			errMsg:  "Mes inválido",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateParams(tt.params)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewIncomesAnalyticsService(t *testing.T) {
	// Crear servicio
	service := NewIncomesAnalyticsService(nil, nil, nil, nil)

	// Verificar
	assert.NotNil(t, service)
}
