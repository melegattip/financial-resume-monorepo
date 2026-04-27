package domain

import "time"

// MonthlyCoachingReport is the AI-generated monthly coaching output for the previous complete month.
type MonthlyCoachingReport struct {
	Month        string          `json:"month"`        // "YYYY-MM" e.g. "2026-02"
	Sentiment    string          `json:"sentiment"`    // "positivo|neutral|desafiante"
	Summary      string          `json:"summary"`
	Wins         []CoachingPoint `json:"wins"`
	Improvements []CoachingPoint `json:"improvements"`
	Actions      []CoachingAction `json:"actions"`
	BehaviorNote string          `json:"behavior_note"`
	GeneratedAt  time.Time       `json:"generated_at"`
}

// CoachingPoint is a single win or improvement item in the monthly coaching report.
type CoachingPoint struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// CoachingAction is a concrete action the user can take in the current month,
// with a deep-link to the relevant section of the app.
type CoachingAction struct {
	Title    string `json:"title"`
	Detail   string `json:"detail"`
	DeepLink string `json:"deep_link"` // e.g. "/budgets", "/savings-goals"
}

// MonthlyCoachingRequest is the body sent by the frontend to POST /ai/monthly-coaching.
type MonthlyCoachingRequest struct {
	FinancialData FinancialAnalysisData `json:"financial_data"`
	PreviousMonth string                `json:"previous_month"` // "YYYY-MM"
	Force         bool                  `json:"force"`          // bypass cache and regenerate
}
