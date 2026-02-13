package migration

import (
	"time"

	"gorm.io/gorm"
)

// --- Source models (gamification-db schema with Financial Health Score fields) ---

// SrcUserGamification represents a row in gamification-db.user_gamification.
type SrcUserGamification struct {
	ID                     string `gorm:"column:id"`
	UserID                 string `gorm:"column:user_id"`
	FinancialHealthScore   int    `gorm:"column:financial_health_score"`
	EngagementComponent    int    `gorm:"column:engagement_component"`
	HealthComponent        int    `gorm:"column:health_component"`
	CurrentLevel           int    `gorm:"column:current_level"`
	InsightsViewed         int    `gorm:"column:insights_viewed"`
	ActionsCompleted       int    `gorm:"column:actions_completed"`
	AchievementsCount      int    `gorm:"column:achievements_count"`
	CurrentStreak          int    `gorm:"column:current_streak"`
	LastActivity           *time.Time
	LastScoreCalculation   *time.Time
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

func (SrcUserGamification) TableName() string { return "user_gamification" }

// SrcAchievement represents a row in gamification-db.achievements.
type SrcAchievement struct {
	ID          string `gorm:"column:id"`
	UserID      string `gorm:"column:user_id"`
	Type        string `gorm:"column:type"`
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
	Points      int    `gorm:"column:points"`
	Progress    int    `gorm:"column:progress"`
	Target      int    `gorm:"column:target"`
	Completed   bool   `gorm:"column:completed"`
	UnlockedAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (SrcAchievement) TableName() string { return "achievements" }

// SrcUserAction represents a row in gamification-db.user_actions.
type SrcUserAction struct {
	ID                string `gorm:"column:id"`
	UserID            string `gorm:"column:user_id"`
	ActionType        string `gorm:"column:action_type"`
	EntityType        string `gorm:"column:entity_type"`
	EntityID          string `gorm:"column:entity_id"`
	ScoreContribution int    `gorm:"column:score_contribution"`
	ComponentAffected string `gorm:"column:component_affected"`
	Description       string `gorm:"column:description"`
	CreatedAt         time.Time
}

func (SrcUserAction) TableName() string { return "user_actions" }

// SrcChallenge represents a row in gamification-db.challenges.
type SrcChallenge struct {
	ID               string `gorm:"column:id"`
	ChallengeKey     string `gorm:"column:challenge_key"`
	Name             string `gorm:"column:name"`
	Description      string `gorm:"column:description"`
	ChallengeType    string `gorm:"column:challenge_type"`
	Icon             string `gorm:"column:icon"`
	XPReward         int    `gorm:"column:xp_reward"`
	RequirementType  string `gorm:"column:requirement_type"`
	RequirementTarget int   `gorm:"column:requirement_target"`
	RequirementData  string `gorm:"column:requirement_data"` // JSONB stored as string
	Active           bool   `gorm:"column:active"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (SrcChallenge) TableName() string { return "challenges" }

// SrcUserChallenge represents a row in gamification-db.user_challenges.
type SrcUserChallenge struct {
	ID            string `gorm:"column:id"`
	UserID        string `gorm:"column:user_id"`
	ChallengeID   string `gorm:"column:challenge_id"`
	ChallengeDate time.Time `gorm:"column:challenge_date"`
	Progress      int    `gorm:"column:progress"`
	Target        int    `gorm:"column:target"`
	Completed     bool   `gorm:"column:completed"`
	CompletedAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (SrcUserChallenge) TableName() string { return "user_challenges" }

// SrcChallengeProgressTracking represents a row in gamification-db.challenge_progress_tracking.
type SrcChallengeProgressTracking struct {
	ID             string `gorm:"column:id"`
	UserID         string `gorm:"column:user_id"`
	ChallengeDate  time.Time `gorm:"column:challenge_date"`
	ActionType     string `gorm:"column:action_type"`
	EntityType     string `gorm:"column:entity_type"`
	Count          int    `gorm:"column:count"`
	UniqueEntities string `gorm:"column:unique_entities"` // JSONB stored as string
	CreatedAt      time.Time
}

func (SrcChallengeProgressTracking) TableName() string { return "challenge_progress_tracking" }

// --- Target models for gamification tables (used for AutoMigrate) ---

// TgtUserGamification is the target schema for user_gamification.
type TgtUserGamification struct {
	ID                     string `gorm:"type:varchar(255);primaryKey"`
	UserID                 string `gorm:"type:varchar(255);uniqueIndex;not null"`
	FinancialHealthScore   int    `gorm:"default:1"`
	EngagementComponent    int    `gorm:"default:0"`
	HealthComponent        int    `gorm:"default:0"`
	CurrentLevel           int    `gorm:"default:1"`
	InsightsViewed         int    `gorm:"default:0"`
	ActionsCompleted       int    `gorm:"default:0"`
	AchievementsCount      int    `gorm:"default:0"`
	CurrentStreak          int    `gorm:"default:0"`
	LastActivity           *time.Time
	LastScoreCalculation   *time.Time
	CreatedAt              time.Time
	UpdatedAt              time.Time
	DeletedAt              gorm.DeletedAt `gorm:"index"`
}

func (TgtUserGamification) TableName() string { return "user_gamification" }

// TgtAchievement is the target schema for achievements.
type TgtAchievement struct {
	ID          string `gorm:"type:varchar(255);primaryKey"`
	UserID      string `gorm:"type:varchar(255);not null;index"`
	Type        string `gorm:"type:varchar(100);not null"`
	Name        string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text"`
	Points      int    `gorm:"default:0"`
	Progress    int    `gorm:"default:0"`
	Target      int    `gorm:"not null"`
	Completed   bool   `gorm:"default:false"`
	UnlockedAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (TgtAchievement) TableName() string { return "achievements" }

// TgtUserAction is the target schema for user_actions.
type TgtUserAction struct {
	ID                string `gorm:"type:varchar(255);primaryKey"`
	UserID            string `gorm:"type:varchar(255);not null;index"`
	ActionType        string `gorm:"type:varchar(100);not null"`
	EntityType        string `gorm:"type:varchar(100);not null"`
	EntityID          string `gorm:"type:varchar(255)"`
	ScoreContribution int    `gorm:"default:0"`
	ComponentAffected string `gorm:"type:varchar(50)"`
	Description       string `gorm:"type:text"`
	CreatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

func (TgtUserAction) TableName() string { return "user_actions" }

// TgtChallenge is the target schema for challenges.
type TgtChallenge struct {
	ID                string `gorm:"type:varchar(255);primaryKey"`
	ChallengeKey      string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Name              string `gorm:"type:varchar(255);not null"`
	Description       string `gorm:"type:text;not null"`
	ChallengeType     string `gorm:"type:varchar(50);not null"`
	Icon              string `gorm:"type:varchar(20);default:'🎯'"`
	XPReward          int    `gorm:"default:0"`
	RequirementType   string `gorm:"type:varchar(100);not null"`
	RequirementTarget int    `gorm:"not null"`
	RequirementData   string `gorm:"type:jsonb"`
	Active            bool   `gorm:"default:true"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

func (TgtChallenge) TableName() string { return "challenges" }

// TgtUserChallenge is the target schema for user_challenges.
type TgtUserChallenge struct {
	ID            string    `gorm:"type:varchar(255);primaryKey"`
	UserID        string    `gorm:"type:varchar(255);not null;index"`
	ChallengeID   string    `gorm:"type:varchar(255);not null;index"`
	ChallengeDate time.Time `gorm:"type:date;not null"`
	Progress      int       `gorm:"default:0"`
	Target        int       `gorm:"not null"`
	Completed     bool      `gorm:"default:false"`
	CompletedAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func (TgtUserChallenge) TableName() string { return "user_challenges" }

// TgtChallengeProgressTracking is the target schema for challenge_progress_tracking.
type TgtChallengeProgressTracking struct {
	ID             string    `gorm:"type:varchar(255);primaryKey"`
	UserID         string    `gorm:"type:varchar(255);not null;index"`
	ChallengeDate  time.Time `gorm:"type:date;not null"`
	ActionType     string    `gorm:"type:varchar(100);not null"`
	EntityType     string    `gorm:"type:varchar(100)"`
	Count          int       `gorm:"default:1"`
	UniqueEntities string    `gorm:"type:jsonb"`
	CreatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func (TgtChallengeProgressTracking) TableName() string { return "challenge_progress_tracking" }
