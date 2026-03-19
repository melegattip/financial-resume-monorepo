package domain

import "time"

// EducationCard is a single AI-generated financial education card,
// personalized to the user's current financial situation.
type EducationCard struct {
	Topic      string `json:"topic"`       // "emergencia|presupuesto|deuda|ahorro|inversión|impuestos"
	Title      string `json:"title"`
	Summary    string `json:"summary"`     // 2-3 sentences explaining the concept
	KeyConcept string `json:"key_concept"` // memorable callout phrase
	CTA        string `json:"cta"`         // button label, max 35 chars
	DeepLink   string `json:"deep_link"`   // e.g. "/savings-goals", "/budgets"
	Difficulty string `json:"difficulty"`  // "básico|intermedio|avanzado"
}

// EducationContent wraps the slice of cards returned by the AI service.
type EducationContent struct {
	Cards       []EducationCard `json:"cards"`
	GeneratedAt time.Time       `json:"generated_at"`
}

// EducationRequest is the body sent by the frontend to POST /ai/education-cards.
type EducationRequest struct {
	FinancialData FinancialAnalysisData `json:"financial_data"`
}
