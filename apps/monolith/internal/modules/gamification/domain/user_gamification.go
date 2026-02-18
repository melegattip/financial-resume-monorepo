package domain

import "time"

// xpThresholds defines the cumulative XP required to reach each level (index = level - 1).
var xpThresholds = []int{0, 75, 200, 400, 700, 1200, 1800, 2600, 3600, 5500}

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
	ID                  string
	UserID              string
	TotalXP             int
	CurrentLevel        int
	InsightsViewed      int
	ActionsCompleted    int
	AchievementsCount   int
	CurrentStreak       int
	LastActivity        time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time
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
