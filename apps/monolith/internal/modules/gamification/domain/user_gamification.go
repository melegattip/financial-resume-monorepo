package domain

import "time"

// xpThresholds defines the cumulative points required to reach each level (index = level - 1).
// The scale runs 0–1000: Level 1 starts at 0, Level 10 (max) is reached at 1000 points.
var xpThresholds = []int{0, 50, 100, 175, 275, 400, 550, 700, 850, 1000}

// maxScore is the maximum score a user can display (Level 10 cap).
const maxScore = 1000

// levelNames maps level number (1-based) to its display name.
var levelNames = []string{
	"Financial Newbie",
	"Money Tracker",
	"Smart Saver",
	"Budget Master",
	"Financial Planner",
	"Investment Seeker",
	"Wealth Builder",
	"Financial Strategist",
	"Money Mentor",
	"Financial Magnate",
}

// UserGamification is the aggregate root for a user's gamification state.
type UserGamification struct {
	ID                  string     `json:"id"`
	UserID              string     `json:"user_id"`
	TotalXP             int        `json:"total_xp"`
	CurrentLevel        int        `json:"current_level"`
	InsightsViewed      int        `json:"insights_viewed"`
	ActionsCompleted    int        `json:"actions_completed"`
	AchievementsCount   int        `json:"achievements_count"`
	CurrentStreak       int        `json:"current_streak"`
	LastActivity        time.Time  `json:"last_activity"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty"`
}

// NewUserGamification creates a new UserGamification for a user, starting at level 1.
func NewUserGamification(userID string) *UserGamification {
	now := time.Now().UTC()
	return &UserGamification{
		UserID:       userID,
		TotalXP:      0,
		CurrentLevel: 1,
		LastActivity: now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// CalculateLevel derives the correct level from TotalXP using the defined thresholds.
// Returns a value between 1 and 10.
func (g *UserGamification) CalculateLevel() int {
	level := 1
	for i, threshold := range xpThresholds {
		if g.TotalXP >= threshold {
			level = i + 1
		}
	}
	if level > 10 {
		level = 10
	}
	return level
}

// XPToNextLevel returns the amount of XP still needed to reach the next level.
// Returns 0 when the user is already at the maximum level (10).
func (g *UserGamification) XPToNextLevel() int {
	if g.CurrentLevel >= 10 {
		return 0
	}
	nextThreshold := xpThresholds[g.CurrentLevel] // index = next level - 1
	remaining := nextThreshold - g.TotalXP
	if remaining < 0 {
		return 0
	}
	return remaining
}

// ProgressToNextLevel returns an integer percentage (0–100) indicating how far the
// user has progressed toward the next level.
func (g *UserGamification) ProgressToNextLevel() int {
	if g.CurrentLevel >= 10 {
		return 100
	}
	currentThreshold := xpThresholds[g.CurrentLevel-1]
	nextThreshold := xpThresholds[g.CurrentLevel]
	levelRange := nextThreshold - currentThreshold
	if levelRange <= 0 {
		return 100
	}
	xpInLevel := g.TotalXP - currentThreshold
	if xpInLevel < 0 {
		return 0
	}
	pct := (xpInLevel * 100) / levelRange
	if pct > 100 {
		return 100
	}
	return pct
}

// GetLevelName returns the human-readable name for the user's current level.
func (g *UserGamification) GetLevelName() string {
	idx := g.CurrentLevel - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= len(levelNames) {
		idx = len(levelNames) - 1
	}
	return levelNames[idx]
}

// Score returns the display score clamped to the 0–1000 range.
// Users who accumulated XP before the rescale may have TotalXP > 1000;
// this method ensures the displayed value never exceeds the max.
func (g *UserGamification) Score() int {
	if g.TotalXP > maxScore {
		return maxScore
	}
	return g.TotalXP
}
