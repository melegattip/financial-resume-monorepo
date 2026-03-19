package email

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNoOpEmailService_SendMonthlyCoachingReport(t *testing.T) {
	svc := &NoOpEmailService{logger: zerolog.Nop()}
	data := CoachingReportEmailData{
		Month:     "2026-02",
		Sentiment: "positivo",
		Summary:   "Buen mes",
		Wins:      []CoachingEmailPoint{{Title: "Win", Description: "Desc"}},
	}
	err := svc.SendMonthlyCoachingReport("test@test.com", "Juan", data)
	assert.NoError(t, err)
}

func TestNoOpEmailService_SendMonthlyCoachingReport_EmptyAddress(t *testing.T) {
	svc := &NoOpEmailService{logger: zerolog.Nop()}
	err := svc.SendMonthlyCoachingReport("", "", CoachingReportEmailData{Month: "2026-02"})
	assert.NoError(t, err)
}

func TestBuildMonthlyCoachingHTML_ContainsMonth(t *testing.T) {
	data := CoachingReportEmailData{
		Month:     "2026-02",
		Sentiment: "positivo",
		Summary:   "Excelente mes",
		Wins: []CoachingEmailPoint{
			{Title: "Ahorro", Description: "Guardaste $500"},
		},
		Improvements: []CoachingEmailPoint{
			{Title: "Delivery", Description: "Reducí gastos"},
		},
		Actions: []CoachingEmailAction{
			{Title: "Crear presupuesto", Detail: "Para delivery"},
		},
	}
	html := buildMonthlyCoachingHTML("Juan", data)
	assert.Contains(t, html, "2026-02")
	assert.Contains(t, html, "Excelente mes")
	assert.Contains(t, html, "Ahorro")
	assert.Contains(t, html, "Delivery")
	assert.Contains(t, html, "Crear presupuesto")
}

func TestBuildMonthlyCoachingHTML_Sentiment_Desafiante(t *testing.T) {
	data := CoachingReportEmailData{Sentiment: "desafiante", Month: "2026-02"}
	html := buildMonthlyCoachingHTML("", data)
	assert.Contains(t, html, "Mes Desafiante")
}

func TestBuildMonthlyCoachingHTML_Sentiment_Neutral(t *testing.T) {
	data := CoachingReportEmailData{Sentiment: "neutral", Month: "2026-02"}
	html := buildMonthlyCoachingHTML("", data)
	assert.Contains(t, html, "Mes Neutral")
}

func TestBuildMonthlyCoachingHTML_Sentiment_Positivo(t *testing.T) {
	data := CoachingReportEmailData{Sentiment: "positivo", Month: "2026-02"}
	html := buildMonthlyCoachingHTML("", data)
	assert.Contains(t, html, "Mes Positivo")
}

func TestBuildMonthlyCoachingHTML_WithBehaviorNote(t *testing.T) {
	data := CoachingReportEmailData{
		Month:        "2026-02",
		BehaviorNote: "Tu consistencia mejoró",
	}
	html := buildMonthlyCoachingHTML("Ana", data)
	assert.Contains(t, html, "Tu consistencia mejoró")
}
