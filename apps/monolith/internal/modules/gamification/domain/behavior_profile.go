package domain

import (
	"math"
	"time"
)

// BehaviorProfile aggregates a user's financial behavioral metrics
// derived from user_actions counts. No new DB columns are needed.
type BehaviorProfile struct {
	UserID    string `json:"user_id"`
	LevelName string `json:"level_name"`
	TotalXP   int    `json:"total_xp"`

	CurrentLevel  int `json:"current_level"`
	CurrentStreak int `json:"current_streak"`
	DaysActive    int `json:"days_active"` // days since gamification was initialized

	AchievementsCompleted int `json:"achievements_completed"`

	// Derived from user_actions COUNT by type (no new DB columns).
	BudgetsCreated           int `json:"budgets_created"`
	BudgetComplianceEvents   int `json:"budget_compliance_events"` // stay_within_budget
	SavingsGoalsCreated      int `json:"savings_goals_created"`
	SavingsDeposits          int `json:"savings_deposits"`
	SavingsGoalsAchieved     int `json:"savings_goals_achieved"`
	RecurringSetups          int `json:"recurring_setups"`
	AIRecommendationsApplied int `json:"ai_recommendations_applied"`
	AnalyticsViewsCount      int `json:"analytics_views_count"`

	// Pre-computed dimension scores (0-100).
	ConsistencyScore int `json:"consistency_score"` // streak + tenure
	DisciplineScore  int `json:"discipline_score"`  // planning actions
	EngagementScore  int `json:"engagement_score"`  // AI + analytics use

	ComputedAt time.Time `json:"computed_at"`
}

// ComputeDimensionScores fills the three pre-computed dimension scores.
// Call this after all action counts are set.
func (b *BehaviorProfile) ComputeDimensionScores() {
	// Consistency: how regular and long-tenured the user is.
	streakFactor := math.Min(float64(b.CurrentStreak)/30.0, 1.0) * 60
	tenureFactor := math.Min(float64(b.DaysActive)/90.0, 1.0) * 40
	b.ConsistencyScore = int(streakFactor + tenureFactor)

	// Discipline: how much the user plans and executes savings/budgets.
	discipline := b.BudgetsCreated*25 + b.SavingsGoalsCreated*25 +
		b.RecurringSetups*20 + b.BudgetComplianceEvents*30
	if discipline > 100 {
		discipline = 100
	}
	b.DisciplineScore = discipline

	// Engagement: how actively the user uses analytical and AI features.
	engagement := b.AIRecommendationsApplied*20 + b.RecurringSetups*15 + b.AnalyticsViewsCount*3
	if engagement > 100 {
		engagement = 100
	}
	b.EngagementScore = engagement
}
