// Package domain defines the core business entities and their behavior
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrInvalidGoalAmount is returned when attempting to create a goal with invalid amount
	ErrInvalidGoalAmount = errors.New("goal amount must be positive")
	// ErrInvalidTargetDate is returned when attempting to create a goal with past target date
	ErrInvalidTargetDate = errors.New("target date must be in the future")
	// ErrGoalAlreadyAchieved is returned when attempting to modify an achieved goal
	ErrGoalAlreadyAchieved = errors.New("goal is already achieved")
)

// SavingsGoalCategory represents the category of savings goal
type SavingsGoalCategory string

const (
	SavingsGoalCategoryVacation   SavingsGoalCategory = "vacation"
	SavingsGoalCategoryEmergency  SavingsGoalCategory = "emergency"
	SavingsGoalCategoryHouse      SavingsGoalCategory = "house"
	SavingsGoalCategoryCar        SavingsGoalCategory = "car"
	SavingsGoalCategoryEducation  SavingsGoalCategory = "education"
	SavingsGoalCategoryRetirement SavingsGoalCategory = "retirement"
	SavingsGoalCategoryInvestment SavingsGoalCategory = "investment"
	SavingsGoalCategoryOther      SavingsGoalCategory = "other"
)

// SavingsGoalStatus represents the current status of a savings goal
type SavingsGoalStatus string

const (
	SavingsGoalStatusActive    SavingsGoalStatus = "active"    // Goal is active and in progress
	SavingsGoalStatusAchieved  SavingsGoalStatus = "achieved"  // Goal has been completed
	SavingsGoalStatusPaused    SavingsGoalStatus = "paused"    // Goal is temporarily paused
	SavingsGoalStatusCancelled SavingsGoalStatus = "cancelled" // Goal has been cancelled
)

// SavingsGoalPriority represents the priority level of the goal
type SavingsGoalPriority string

const (
	SavingsGoalPriorityHigh   SavingsGoalPriority = "high"
	SavingsGoalPriorityMedium SavingsGoalPriority = "medium"
	SavingsGoalPriorityLow    SavingsGoalPriority = "low"
)

// SavingsGoal represents a financial savings goal
type SavingsGoal struct {
	ID                string              `json:"id"`
	UserID            string              `json:"user_id"`
	Name              string              `json:"name"`
	Description       string              `json:"description,omitempty"`
	TargetAmount      float64             `json:"target_amount"`
	CurrentAmount     float64             `json:"current_amount"`
	Category          SavingsGoalCategory `json:"category"`
	Priority          SavingsGoalPriority `json:"priority"`
	TargetDate        time.Time           `json:"target_date"`
	Status            SavingsGoalStatus   `json:"status"`
	MonthlyTarget     float64             `json:"monthly_target"`      // Calculated monthly savings needed
	WeeklyTarget      float64             `json:"weekly_target"`       // Calculated weekly savings needed
	DailyTarget       float64             `json:"daily_target"`        // Calculated daily savings needed
	Progress          float64             `json:"progress"`            // Calculated progress percentage (0-1)
	RemainingAmount   float64             `json:"remaining_amount"`    // Calculated remaining amount
	DaysRemaining     int                 `json:"days_remaining"`      // Calculated days remaining
	IsAutoSave        bool                `json:"is_auto_save"`        // Whether to auto-save towards this goal
	AutoSaveAmount    float64             `json:"auto_save_amount"`    // Amount to auto-save periodically
	AutoSaveFrequency string              `json:"auto_save_frequency"` // daily, weekly, monthly
	ImageURL          string              `json:"image_url,omitempty"` // Optional image for motivation
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
	AchievedAt        *time.Time          `json:"achieved_at,omitempty"`
}

// SavingsGoalBuilder implements the Builder pattern for creating savings goals
type SavingsGoalBuilder struct {
	goal *SavingsGoal
}

// NewSavingsGoalBuilder creates a new SavingsGoalBuilder instance
func NewSavingsGoalBuilder() *SavingsGoalBuilder {
	return &SavingsGoalBuilder{
		goal: &SavingsGoal{
			ID:            "goal_" + uuid.New().String()[:8],
			CurrentAmount: 0,
			Status:        SavingsGoalStatusActive,
			Priority:      SavingsGoalPriorityMedium,
			Category:      SavingsGoalCategoryOther,
		},
	}
}

// SetUserID sets the user ID (Builder pattern)
func (b *SavingsGoalBuilder) SetUserID(userID string) *SavingsGoalBuilder {
	b.goal.UserID = userID
	return b
}

// SetName sets the goal name (Builder pattern)
func (b *SavingsGoalBuilder) SetName(name string) *SavingsGoalBuilder {
	b.goal.Name = name
	return b
}

// SetDescription sets the goal description (Builder pattern)
func (b *SavingsGoalBuilder) SetDescription(description string) *SavingsGoalBuilder {
	b.goal.Description = description
	return b
}

// SetTargetAmount sets the target amount (Builder pattern)
func (b *SavingsGoalBuilder) SetTargetAmount(amount float64) *SavingsGoalBuilder {
	b.goal.TargetAmount = amount
	return b
}

// SetCategory sets the goal category (Builder pattern)
func (b *SavingsGoalBuilder) SetCategory(category SavingsGoalCategory) *SavingsGoalBuilder {
	b.goal.Category = category
	return b
}

// SetPriority sets the goal priority (Builder pattern)
func (b *SavingsGoalBuilder) SetPriority(priority SavingsGoalPriority) *SavingsGoalBuilder {
	b.goal.Priority = priority
	return b
}

// SetTargetDate sets the target date (Builder pattern)
func (b *SavingsGoalBuilder) SetTargetDate(targetDate time.Time) *SavingsGoalBuilder {
	b.goal.TargetDate = targetDate
	b.calculateTargets()
	return b
}

// SetAutoSave sets auto-save configuration (Builder pattern)
func (b *SavingsGoalBuilder) SetAutoSave(amount float64, frequency string) *SavingsGoalBuilder {
	b.goal.IsAutoSave = true
	b.goal.AutoSaveAmount = amount
	b.goal.AutoSaveFrequency = frequency
	return b
}

// SetImageURL sets the goal image URL (Builder pattern)
func (b *SavingsGoalBuilder) SetImageURL(imageURL string) *SavingsGoalBuilder {
	b.goal.ImageURL = imageURL
	return b
}

// calculateTargets calculates the required daily, weekly, and monthly targets
func (b *SavingsGoalBuilder) calculateTargets() {
	if b.goal.TargetDate.IsZero() || b.goal.TargetAmount <= 0 {
		return
	}

	now := time.Now()
	daysUntilTarget := int(b.goal.TargetDate.Sub(now).Hours() / 24)

	if daysUntilTarget <= 0 {
		return
	}

	remainingAmount := b.goal.TargetAmount - b.goal.CurrentAmount
	if remainingAmount <= 0 {
		return
	}

	// Calculate targets
	b.goal.DailyTarget = remainingAmount / float64(daysUntilTarget)
	b.goal.WeeklyTarget = b.goal.DailyTarget * 7
	b.goal.MonthlyTarget = remainingAmount / (float64(daysUntilTarget) / 30.0)

	// Update calculated fields
	b.goal.updateCalculatedFields()
}

// UpdateCalculatedFields updates all calculated fields for the goal
func (s *SavingsGoal) UpdateCalculatedFields() {
	s.updateCalculatedFields()
}

// updateCalculatedFields is the internal method to update calculated fields
func (s *SavingsGoal) updateCalculatedFields() {
	// Update progress
	if s.TargetAmount > 0 {
		s.Progress = s.CurrentAmount / s.TargetAmount
		if s.Progress > 1.0 {
			s.Progress = 1.0
		}
	}

	// Update remaining amount
	s.RemainingAmount = s.TargetAmount - s.CurrentAmount
	if s.RemainingAmount < 0 {
		s.RemainingAmount = 0
	}

	// Update days remaining
	if !s.TargetDate.IsZero() {
		daysUntilTarget := int(s.TargetDate.Sub(time.Now()).Hours() / 24)
		if daysUntilTarget < 0 {
			s.DaysRemaining = 0
		} else {
			s.DaysRemaining = daysUntilTarget
		}
	}

	// Recalculate targets if needed
	if !s.TargetDate.IsZero() && s.TargetAmount > 0 && s.RemainingAmount > 0 && s.DaysRemaining > 0 {
		s.DailyTarget = s.RemainingAmount / float64(s.DaysRemaining)
		s.WeeklyTarget = s.DailyTarget * 7
		s.MonthlyTarget = s.RemainingAmount / (float64(s.DaysRemaining) / 30.0)
	}
}

// Build creates and validates the savings goal (Builder pattern)
func (b *SavingsGoalBuilder) Build() (*SavingsGoal, error) {
	if err := b.goal.Validate(); err != nil {
		return nil, err
	}

	now := time.Now()
	b.goal.CreatedAt = now
	b.goal.UpdatedAt = now

	return b.goal, nil
}

// Validate validates the savings goal business rules (Single Responsibility Principle)
func (s *SavingsGoal) Validate() error {
	if s.UserID == "" {
		return errors.New("user ID is required")
	}

	if s.Name == "" {
		return errors.New("goal name is required")
	}

	if s.TargetAmount <= 0 {
		return ErrInvalidGoalAmount
	}

	if s.CurrentAmount < 0 {
		return errors.New("current amount cannot be negative")
	}

	if s.CurrentAmount > s.TargetAmount {
		return errors.New("current amount cannot exceed target amount")
	}

	// Allow target dates in the past for existing goals (overdue goals are valid in updates).
	// Creation enforces future target date at the use case level.

	if s.IsAutoSave {
		if s.AutoSaveAmount <= 0 {
			return errors.New("auto-save amount must be positive")
		}

		validFrequencies := map[string]bool{
			"daily":   true,
			"weekly":  true,
			"monthly": true,
		}

		if !validFrequencies[s.AutoSaveFrequency] {
			return errors.New("auto-save frequency must be daily, weekly, or monthly")
		}
	}

	return nil
}

// AddSavings adds money to the goal (Single Responsibility)
func (s *SavingsGoal) AddSavings(amount float64) error {
	if amount <= 0 {
		return errors.New("savings amount must be positive")
	}

	if s.Status == SavingsGoalStatusAchieved {
		return ErrGoalAlreadyAchieved
	}

	if s.Status == SavingsGoalStatusCancelled {
		return errors.New("cannot add savings to cancelled goal")
	}

	s.CurrentAmount += amount

	// Check if goal is achieved
	if s.CurrentAmount >= s.TargetAmount {
		s.Status = SavingsGoalStatusAchieved
		now := time.Now()
		s.AchievedAt = &now
	}

	s.UpdatedAt = time.Now()
	s.updateCalculatedFields()

	return nil
}

// WithdrawSavings removes money from the goal (Single Responsibility)
func (s *SavingsGoal) WithdrawSavings(amount float64) error {
	if amount <= 0 {
		return errors.New("withdrawal amount must be positive")
	}

	if s.CurrentAmount < amount {
		return errors.New("insufficient savings in goal")
	}

	s.CurrentAmount -= amount

	// If goal was achieved but now isn't, change status back to active
	if s.Status == SavingsGoalStatusAchieved && s.CurrentAmount < s.TargetAmount {
		s.Status = SavingsGoalStatusActive
		s.AchievedAt = nil
	}

	s.UpdatedAt = time.Now()
	s.updateCalculatedFields()

	return nil
}

// recalculateTargets recalculates the required savings targets
func (s *SavingsGoal) recalculateTargets() {
	if s.TargetDate.IsZero() || s.Status == SavingsGoalStatusAchieved {
		s.DailyTarget = 0
		s.WeeklyTarget = 0
		s.MonthlyTarget = 0
		return
	}

	now := time.Now()
	daysUntilTarget := int(s.TargetDate.Sub(now).Hours() / 24)

	if daysUntilTarget <= 0 {
		s.DailyTarget = 0
		s.WeeklyTarget = 0
		s.MonthlyTarget = 0
		return
	}

	remainingAmount := s.TargetAmount - s.CurrentAmount
	if remainingAmount <= 0 {
		s.DailyTarget = 0
		s.WeeklyTarget = 0
		s.MonthlyTarget = 0
		return
	}

	s.DailyTarget = remainingAmount / float64(daysUntilTarget)
	s.WeeklyTarget = s.DailyTarget * 7
	s.MonthlyTarget = remainingAmount / (float64(daysUntilTarget) / 30.0)
}

// GetProgress returns the progress percentage (0.0-1.0)
func (s *SavingsGoal) GetProgress() float64 {
	if s.TargetAmount == 0 {
		return 0
	}
	progress := s.CurrentAmount / s.TargetAmount
	if progress > 1.0 {
		return 1.0
	}
	return progress
}

// GetRemainingAmount returns the remaining amount needed
func (s *SavingsGoal) GetRemainingAmount() float64 {
	remaining := s.TargetAmount - s.CurrentAmount
	if remaining < 0 {
		return 0
	}
	return remaining
}

// GetDaysRemaining returns the number of days until target date
func (s *SavingsGoal) GetDaysRemaining() int {
	if s.TargetDate.IsZero() {
		return 0
	}

	days := int(s.TargetDate.Sub(time.Now()).Hours() / 24)
	if days < 0 {
		return 0
	}
	return days
}

// IsOverdue returns true if the target date has passed and goal is not achieved
func (s *SavingsGoal) IsOverdue() bool {
	if s.TargetDate.IsZero() || s.Status == SavingsGoalStatusAchieved {
		return false
	}

	return time.Now().After(s.TargetDate)
}

// IsOnTrack returns true if the current progress is on track to meet the target date
func (s *SavingsGoal) IsOnTrack() bool {
	if s.TargetDate.IsZero() || s.Status == SavingsGoalStatusAchieved {
		return true
	}

	now := time.Now()
	totalDays := s.TargetDate.Sub(s.CreatedAt).Hours() / 24
	daysPassed := now.Sub(s.CreatedAt).Hours() / 24

	if totalDays <= 0 {
		return true
	}

	expectedProgress := daysPassed / totalDays
	actualProgress := s.GetProgress()

	// Consider on track if within 10% of expected progress
	return actualProgress >= (expectedProgress - 0.1)
}

// Pause pauses the goal
func (s *SavingsGoal) Pause() error {
	if s.Status == SavingsGoalStatusAchieved {
		return ErrGoalAlreadyAchieved
	}

	if s.Status == SavingsGoalStatusCancelled {
		return errors.New("cannot pause cancelled goal")
	}

	s.Status = SavingsGoalStatusPaused
	s.UpdatedAt = time.Now()

	return nil
}

// Resume resumes a paused goal
func (s *SavingsGoal) Resume() error {
	if s.Status != SavingsGoalStatusPaused {
		return errors.New("can only resume paused goals")
	}

	s.Status = SavingsGoalStatusActive
	s.UpdatedAt = time.Now()
	s.recalculateTargets()

	return nil
}

// Cancel cancels the goal
func (s *SavingsGoal) Cancel() error {
	if s.Status == SavingsGoalStatusAchieved {
		return ErrGoalAlreadyAchieved
	}

	s.Status = SavingsGoalStatusCancelled
	s.UpdatedAt = time.Now()

	return nil
}

// SavingsGoalSummary provides aggregated savings goal information
type SavingsGoalSummary struct {
	TotalGoals      int     `json:"total_goals"`
	ActiveGoals     int     `json:"active_goals"`
	AchievedGoals   int     `json:"achieved_goals"`
	PausedGoals     int     `json:"paused_goals"`
	CancelledGoals  int     `json:"cancelled_goals"`
	TotalTarget     float64 `json:"total_target"`
	TotalSaved      float64 `json:"total_saved"`
	TotalRemaining  float64 `json:"total_remaining"`
	AverageProgress float64 `json:"average_progress"`
	OverdueGoals    int     `json:"overdue_goals"`
	OnTrackGoals    int     `json:"on_track_goals"`
}

// SavingsTransactionType represents the type of savings transaction
type SavingsTransactionType string

const (
	SavingsTransactionTypeDeposit    SavingsTransactionType = "deposit"
	SavingsTransactionTypeWithdrawal SavingsTransactionType = "withdrawal"
)

// SavingsTransaction represents a transaction towards a savings goal
type SavingsTransaction struct {
	ID          string                 `json:"id"`
	GoalID      string                 `json:"goal_id"`
	UserID      string                 `json:"user_id"`
	Amount      float64                `json:"amount"`
	Type        SavingsTransactionType `json:"type"`
	Description string                 `json:"description"`
	CreatedAt   time.Time              `json:"created_at"`
}

// NewSavingsGoalID generates a new savings goal ID
func NewSavingsGoalID() string {
	return "goal_" + uuid.New().String()[:8]
}

// NewSavingsTransactionID generates a new savings transaction ID
func NewSavingsTransactionID() string {
	return "stxn_" + uuid.New().String()[:8]
}
