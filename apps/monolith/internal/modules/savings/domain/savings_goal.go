package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrInvalidGoalAmount is returned when attempting to create a goal with an invalid amount
	ErrInvalidGoalAmount = errors.New("goal amount must be positive")
	// ErrInvalidTargetDate is returned when attempting to create a goal with a past target date
	ErrInvalidTargetDate = errors.New("target date must be in the future")
	// ErrGoalAlreadyAchieved is returned when attempting to modify an achieved goal
	ErrGoalAlreadyAchieved = errors.New("goal is already achieved")
)

// SavingsGoalCategory represents the category of a savings goal
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
	SavingsGoalStatusActive    SavingsGoalStatus = "active"
	SavingsGoalStatusAchieved  SavingsGoalStatus = "achieved"
	SavingsGoalStatusPaused    SavingsGoalStatus = "paused"
	SavingsGoalStatusCancelled SavingsGoalStatus = "cancelled"
)

// SavingsGoalPriority represents the priority level of a savings goal
type SavingsGoalPriority string

const (
	SavingsGoalPriorityHigh   SavingsGoalPriority = "high"
	SavingsGoalPriorityMedium SavingsGoalPriority = "medium"
	SavingsGoalPriorityLow    SavingsGoalPriority = "low"
)

// SavingsTransactionType represents the type of a savings transaction
type SavingsTransactionType string

const (
	SavingsTransactionTypeDeposit    SavingsTransactionType = "deposit"
	SavingsTransactionTypeWithdrawal SavingsTransactionType = "withdrawal"
)

// SavingsGoal represents a financial savings goal
type SavingsGoal struct {
	ID                string
	UserID            string
	TenantID          string
	Name              string
	Description       string
	TargetAmount      float64
	CurrentAmount     float64
	Category          SavingsGoalCategory
	Priority          SavingsGoalPriority
	TargetDate        time.Time
	Status            SavingsGoalStatus
	MonthlyTarget     float64 // Calculated monthly savings needed
	WeeklyTarget      float64 // Calculated weekly savings needed
	DailyTarget       float64 // Calculated daily savings needed
	Progress          float64 // 0.0-1.0
	RemainingAmount   float64
	DaysRemaining     int
	IsAutoSave        bool
	AutoSaveAmount    float64
	AutoSaveFrequency string // "daily", "weekly", "monthly"
	ImageURL          string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	AchievedAt        *time.Time
	DeletedAt         *time.Time
}

// SavingsTransaction represents a deposit or withdrawal towards a savings goal
type SavingsTransaction struct {
	ID          string
	GoalID      string
	UserID      string
	Amount      float64
	Type        SavingsTransactionType
	Description string
	CreatedAt   time.Time
}

// NewSavingsGoalID generates a new savings goal ID
func NewSavingsGoalID() string {
	return "goal_" + uuid.New().String()[:8]
}

// NewSavingsTransactionID generates a new savings transaction ID
func NewSavingsTransactionID() string {
	return "stxn_" + uuid.New().String()[:8]
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

// Validate validates the savings goal business rules
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

// UpdateCalculatedFields recalculates MonthlyTarget, WeeklyTarget, DailyTarget, Progress, RemainingAmount, DaysRemaining
func (s *SavingsGoal) UpdateCalculatedFields() {
	// Progress
	if s.TargetAmount > 0 {
		s.Progress = s.CurrentAmount / s.TargetAmount
		if s.Progress > 1.0 {
			s.Progress = 1.0
		}
	}

	// Remaining amount
	s.RemainingAmount = s.TargetAmount - s.CurrentAmount
	if s.RemainingAmount < 0 {
		s.RemainingAmount = 0
	}

	// Days remaining
	if !s.TargetDate.IsZero() {
		days := int(time.Until(s.TargetDate).Hours() / 24)
		if days < 0 {
			s.DaysRemaining = 0
		} else {
			s.DaysRemaining = days
		}
	}

	// Recalculate periodic targets
	if !s.TargetDate.IsZero() && s.TargetAmount > 0 && s.RemainingAmount > 0 && s.DaysRemaining > 0 {
		s.DailyTarget = s.RemainingAmount / float64(s.DaysRemaining)
		s.WeeklyTarget = s.DailyTarget * 7
		s.MonthlyTarget = s.RemainingAmount / (float64(s.DaysRemaining) / 30.0)
	} else {
		s.DailyTarget = 0
		s.WeeklyTarget = 0
		s.MonthlyTarget = 0
	}
}

// AddSavings adds money to the goal and marks it as achieved if it reaches the target
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

	if s.CurrentAmount >= s.TargetAmount {
		s.Status = SavingsGoalStatusAchieved
		now := time.Now().UTC()
		s.AchievedAt = &now
	}

	s.UpdatedAt = time.Now().UTC()
	s.UpdateCalculatedFields()
	return nil
}

// WithdrawSavings removes money from the goal
func (s *SavingsGoal) WithdrawSavings(amount float64) error {
	if amount <= 0 {
		return errors.New("withdrawal amount must be positive")
	}
	if s.CurrentAmount < amount {
		return errors.New("insufficient savings in goal")
	}

	s.CurrentAmount -= amount

	// If goal was achieved but is no longer, change status back to active
	if s.Status == SavingsGoalStatusAchieved && s.CurrentAmount < s.TargetAmount {
		s.Status = SavingsGoalStatusActive
		s.AchievedAt = nil
	}

	s.UpdatedAt = time.Now().UTC()
	s.UpdateCalculatedFields()
	return nil
}

// Pause pauses an active goal
func (s *SavingsGoal) Pause() error {
	if s.Status == SavingsGoalStatusAchieved {
		return ErrGoalAlreadyAchieved
	}
	if s.Status == SavingsGoalStatusCancelled {
		return errors.New("cannot pause cancelled goal")
	}
	s.Status = SavingsGoalStatusPaused
	s.UpdatedAt = time.Now().UTC()
	return nil
}

// Resume resumes a paused goal
func (s *SavingsGoal) Resume() error {
	if s.Status != SavingsGoalStatusPaused {
		return errors.New("can only resume paused goals")
	}
	s.Status = SavingsGoalStatusActive
	s.UpdatedAt = time.Now().UTC()
	s.UpdateCalculatedFields()
	return nil
}

// Cancel cancels a goal
func (s *SavingsGoal) Cancel() error {
	if s.Status == SavingsGoalStatusAchieved {
		return ErrGoalAlreadyAchieved
	}
	s.Status = SavingsGoalStatusCancelled
	s.UpdatedAt = time.Now().UTC()
	return nil
}

// GetProgress returns the progress as a value between 0.0 and 1.0
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

// GetRemainingAmount returns the amount still needed to reach the target
func (s *SavingsGoal) GetRemainingAmount() float64 {
	remaining := s.TargetAmount - s.CurrentAmount
	if remaining < 0 {
		return 0
	}
	return remaining
}

// GetDaysRemaining returns the number of days until the target date
func (s *SavingsGoal) GetDaysRemaining() int {
	if s.TargetDate.IsZero() {
		return 0
	}
	days := int(time.Until(s.TargetDate).Hours() / 24)
	if days < 0 {
		return 0
	}
	return days
}

// IsOverdue returns true if the target date has passed and the goal is not achieved
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

// SavingsGoalBuilder implements the Builder pattern for creating savings goals
type SavingsGoalBuilder struct {
	goal *SavingsGoal
}

// NewSavingsGoalBuilder creates a new SavingsGoalBuilder
func NewSavingsGoalBuilder() *SavingsGoalBuilder {
	return &SavingsGoalBuilder{
		goal: &SavingsGoal{
			ID:            NewSavingsGoalID(),
			CurrentAmount: 0,
			Status:        SavingsGoalStatusActive,
			Priority:      SavingsGoalPriorityMedium,
			Category:      SavingsGoalCategoryOther,
		},
	}
}

// SetUserID sets the user ID
func (b *SavingsGoalBuilder) SetUserID(userID string) *SavingsGoalBuilder {
	b.goal.UserID = userID
	return b
}

// SetName sets the goal name
func (b *SavingsGoalBuilder) SetName(name string) *SavingsGoalBuilder {
	b.goal.Name = name
	return b
}

// SetDescription sets the goal description
func (b *SavingsGoalBuilder) SetDescription(description string) *SavingsGoalBuilder {
	b.goal.Description = description
	return b
}

// SetTargetAmount sets the target amount
func (b *SavingsGoalBuilder) SetTargetAmount(amount float64) *SavingsGoalBuilder {
	b.goal.TargetAmount = amount
	return b
}

// SetCategory sets the goal category
func (b *SavingsGoalBuilder) SetCategory(category SavingsGoalCategory) *SavingsGoalBuilder {
	b.goal.Category = category
	return b
}

// SetPriority sets the goal priority
func (b *SavingsGoalBuilder) SetPriority(priority SavingsGoalPriority) *SavingsGoalBuilder {
	b.goal.Priority = priority
	return b
}

// SetTargetDate sets the target date and triggers recalculation of periodic targets
func (b *SavingsGoalBuilder) SetTargetDate(targetDate time.Time) *SavingsGoalBuilder {
	b.goal.TargetDate = targetDate
	b.goal.UpdateCalculatedFields()
	return b
}

// SetAutoSave configures auto-save behaviour
func (b *SavingsGoalBuilder) SetAutoSave(amount float64, frequency string) *SavingsGoalBuilder {
	b.goal.IsAutoSave = true
	b.goal.AutoSaveAmount = amount
	b.goal.AutoSaveFrequency = frequency
	return b
}

// SetImageURL sets the optional motivational image URL
func (b *SavingsGoalBuilder) SetImageURL(imageURL string) *SavingsGoalBuilder {
	b.goal.ImageURL = imageURL
	return b
}

// Build validates and returns the constructed SavingsGoal
func (b *SavingsGoalBuilder) Build() (*SavingsGoal, error) {
	if err := b.goal.Validate(); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	b.goal.CreatedAt = now
	b.goal.UpdatedAt = now
	b.goal.UpdateCalculatedFields()
	return b.goal, nil
}
