package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// ComputeDimensionScores — ConsistencyScore
// ---------------------------------------------------------------------------

func TestComputeDimensionScores_ConsistencyZero(t *testing.T) {
	b := &BehaviorProfile{CurrentStreak: 0, DaysActive: 0}
	b.ComputeDimensionScores()
	assert.Equal(t, 0, b.ConsistencyScore)
}

func TestComputeDimensionScores_ConsistencyMaxStreak(t *testing.T) {
	// 30-day streak alone gives 60 points; 0 tenure gives 0.
	b := &BehaviorProfile{CurrentStreak: 30, DaysActive: 0}
	b.ComputeDimensionScores()
	assert.Equal(t, 60, b.ConsistencyScore)
}

func TestComputeDimensionScores_ConsistencyMaxTenure(t *testing.T) {
	// 90+ days active gives 40 points; 0 streak gives 0.
	b := &BehaviorProfile{CurrentStreak: 0, DaysActive: 90}
	b.ComputeDimensionScores()
	assert.Equal(t, 40, b.ConsistencyScore)
}

func TestComputeDimensionScores_ConsistencyFull(t *testing.T) {
	// 30-day streak + 90-day tenure = 100.
	b := &BehaviorProfile{CurrentStreak: 30, DaysActive: 90}
	b.ComputeDimensionScores()
	assert.Equal(t, 100, b.ConsistencyScore)
}

func TestComputeDimensionScores_ConsistencyBeyondMax(t *testing.T) {
	// Streak and tenure capped at 1.0 factor — exceeding 30/90 does not overflow.
	b := &BehaviorProfile{CurrentStreak: 999, DaysActive: 999}
	b.ComputeDimensionScores()
	assert.Equal(t, 100, b.ConsistencyScore)
}

func TestComputeDimensionScores_ConsistencyIntermediate(t *testing.T) {
	// 15-day streak → 0.5 * 60 = 30; 45-day tenure → 0.5 * 40 = 20 → total 50.
	b := &BehaviorProfile{CurrentStreak: 15, DaysActive: 45}
	b.ComputeDimensionScores()
	assert.Equal(t, 50, b.ConsistencyScore)
}

// ---------------------------------------------------------------------------
// ComputeDimensionScores — DisciplineScore
// ---------------------------------------------------------------------------

func TestComputeDimensionScores_DisciplineZero(t *testing.T) {
	b := &BehaviorProfile{}
	b.ComputeDimensionScores()
	assert.Equal(t, 0, b.DisciplineScore)
}

func TestComputeDimensionScores_DisciplineBudgetOnly(t *testing.T) {
	// 1 budget = 25 pts.
	b := &BehaviorProfile{BudgetsCreated: 1}
	b.ComputeDimensionScores()
	assert.Equal(t, 25, b.DisciplineScore)
}

func TestComputeDimensionScores_DisciplineSavingsGoalOnly(t *testing.T) {
	// 1 savings goal = 25 pts.
	b := &BehaviorProfile{SavingsGoalsCreated: 1}
	b.ComputeDimensionScores()
	assert.Equal(t, 25, b.DisciplineScore)
}

func TestComputeDimensionScores_DisciplineRecurringOnly(t *testing.T) {
	// 1 recurring = 20 pts.
	b := &BehaviorProfile{RecurringSetups: 1}
	b.ComputeDimensionScores()
	assert.Equal(t, 20, b.DisciplineScore)
}

func TestComputeDimensionScores_DisciplineComplianceOnly(t *testing.T) {
	// 1 compliance event = 30 pts.
	b := &BehaviorProfile{BudgetComplianceEvents: 1}
	b.ComputeDimensionScores()
	assert.Equal(t, 30, b.DisciplineScore)
}

func TestComputeDimensionScores_DisciplineCapped(t *testing.T) {
	// Everything set high — must be capped at 100.
	b := &BehaviorProfile{
		BudgetsCreated:         10,
		SavingsGoalsCreated:    10,
		RecurringSetups:        10,
		BudgetComplianceEvents: 10,
	}
	b.ComputeDimensionScores()
	assert.Equal(t, 100, b.DisciplineScore)
}

func TestComputeDimensionScores_DisciplineExact100(t *testing.T) {
	// 1 budget (25) + 1 savings (25) + 1 recurring (20) + 1 compliance (30) = 100.
	b := &BehaviorProfile{
		BudgetsCreated:         1,
		SavingsGoalsCreated:    1,
		RecurringSetups:        1,
		BudgetComplianceEvents: 1,
	}
	b.ComputeDimensionScores()
	assert.Equal(t, 100, b.DisciplineScore)
}

// ---------------------------------------------------------------------------
// ComputeDimensionScores — EngagementScore
// ---------------------------------------------------------------------------

func TestComputeDimensionScores_EngagementZero(t *testing.T) {
	b := &BehaviorProfile{}
	b.ComputeDimensionScores()
	assert.Equal(t, 0, b.EngagementScore)
}

func TestComputeDimensionScores_EngagementAIOnly(t *testing.T) {
	// 1 AI recommendation = 20 pts.
	b := &BehaviorProfile{AIRecommendationsApplied: 1}
	b.ComputeDimensionScores()
	assert.Equal(t, 20, b.EngagementScore)
}

func TestComputeDimensionScores_EngagementRecurringOnly(t *testing.T) {
	// 1 recurring setup = 15 pts.
	b := &BehaviorProfile{RecurringSetups: 1}
	b.ComputeDimensionScores()
	assert.Equal(t, 15, b.EngagementScore)
}

func TestComputeDimensionScores_EngagementAnalyticsOnly(t *testing.T) {
	// 5 analytics views = 15 pts.
	b := &BehaviorProfile{AnalyticsViewsCount: 5}
	b.ComputeDimensionScores()
	assert.Equal(t, 15, b.EngagementScore)
}

func TestComputeDimensionScores_EngagementCapped(t *testing.T) {
	// High values must be capped at 100.
	b := &BehaviorProfile{
		AIRecommendationsApplied: 10,
		RecurringSetups:          10,
		AnalyticsViewsCount:      100,
	}
	b.ComputeDimensionScores()
	assert.Equal(t, 100, b.EngagementScore)
}

func TestComputeDimensionScores_EngagementIntermediate(t *testing.T) {
	// 2 AI (40) + 2 recurring (30) = 70.
	b := &BehaviorProfile{AIRecommendationsApplied: 2, RecurringSetups: 2}
	b.ComputeDimensionScores()
	assert.Equal(t, 70, b.EngagementScore)
}

// ---------------------------------------------------------------------------
// All three scores computed in a single call
// ---------------------------------------------------------------------------

func TestComputeDimensionScores_FullProfile(t *testing.T) {
	b := &BehaviorProfile{
		CurrentStreak:            15,
		DaysActive:               45,
		BudgetsCreated:           1,
		SavingsGoalsCreated:      1,
		RecurringSetups:          1,
		BudgetComplianceEvents:   0,
		AIRecommendationsApplied: 1,
		AnalyticsViewsCount:      3,
	}
	b.ComputeDimensionScores()

	// Consistency: 15/30*60 + 45/90*40 = 30 + 20 = 50
	assert.Equal(t, 50, b.ConsistencyScore)
	// Discipline: 25 + 25 + 20 + 0 = 70
	assert.Equal(t, 70, b.DisciplineScore)
	// Engagement: 20 + 15 + 9 = 44
	assert.Equal(t, 44, b.EngagementScore)
}
