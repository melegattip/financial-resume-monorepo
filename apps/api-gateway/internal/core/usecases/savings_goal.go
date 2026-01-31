package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/ports"
)

// SavingsGoalUseCase implements the business logic for savings goals (Clean Architecture)
type SavingsGoalUseCase struct {
	repository      ports.SavingsGoalRepository
	notificationSvc ports.SavingsGoalNotificationService
}

// NewSavingsGoalUseCase creates a new instance of SavingsGoalUseCase (Factory pattern)
func NewSavingsGoalUseCase(
	repository ports.SavingsGoalRepository,
	notificationSvc ports.SavingsGoalNotificationService,
) ports.SavingsGoalUseCase {
	return &SavingsGoalUseCase{
		repository:      repository,
		notificationSvc: notificationSvc,
	}
}

// CreateGoal creates a new savings goal with validation (Single Responsibility)
func (uc *SavingsGoalUseCase) CreateGoal(ctx context.Context, request ports.CreateSavingsGoalRequest) (*ports.CreateSavingsGoalResponse, error) {
	// Validate business rules
	if err := uc.validateCreateGoalRequest(request); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Create goal using Builder pattern
	builder := domain.NewSavingsGoalBuilder().
		SetUserID(request.UserID).
		SetName(request.Name).
		SetDescription(request.Description).
		SetTargetAmount(request.TargetAmount).
		SetCategory(request.Category).
		SetPriority(request.Priority).
		SetImageURL(request.ImageURL)

	// Set target date if provided
	if !request.TargetDate.IsZero() {
		builder = builder.SetTargetDate(request.TargetDate)
	}

	// Set auto-save configuration
	if request.IsAutoSave {
		builder = builder.SetAutoSave(request.AutoSaveAmount, request.AutoSaveFrequency)
	}

	goal, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build goal: %w", err)
	}

	// Save to repository
	if err := uc.repository.Create(ctx, goal); err != nil {
		return nil, fmt.Errorf("failed to create goal: %w", err)
	}

	// Convert to response
	response := uc.convertGoalToCreateResponse(goal)
	return response, nil
}

// GetGoal retrieves a specific savings goal (Open/Closed Principle)
func (uc *SavingsGoalUseCase) GetGoal(ctx context.Context, request ports.GetSavingsGoalRequest) (*ports.GetSavingsGoalResponse, error) {
	goal, err := uc.repository.GetByID(ctx, request.UserID, request.GoalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goal: %w", err)
	}

	if goal == nil {
		return nil, errors.New("goal not found")
	}

	response := uc.convertGoalToGetResponse(goal)
	return response, nil
}

// ListGoals retrieves all savings goals for a user with optional filters
func (uc *SavingsGoalUseCase) ListGoals(ctx context.Context, request ports.ListSavingsGoalsRequest) (*ports.ListSavingsGoalsResponse, error) {
	var goals []*domain.SavingsGoal
	var err error

	// Apply filters (Strategy pattern)
	if request.Status != "" {
		goals, err = uc.repository.ListByStatus(ctx, request.UserID, request.Status)
	} else if request.Category != "" {
		goals, err = uc.repository.ListByCategory(ctx, request.UserID, request.Category)
	} else {
		goals, err = uc.repository.List(ctx, request.UserID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list goals: %w", err)
	}

	// Convert to response with summary
	response := &ports.ListSavingsGoalsResponse{
		Goals:   make([]ports.GetSavingsGoalResponse, len(goals)),
		Summary: uc.calculateSummary(goals),
		Count:   len(goals),
	}

	for i, goal := range goals {
		response.Goals[i] = *uc.convertGoalToGetResponse(goal)
	}

	return response, nil
}

// UpdateGoal updates an existing savings goal (Single Responsibility)
func (uc *SavingsGoalUseCase) UpdateGoal(ctx context.Context, request ports.UpdateSavingsGoalRequest) (*ports.UpdateSavingsGoalResponse, error) {
	// Get existing goal
	goal, err := uc.repository.GetByID(ctx, request.UserID, request.GoalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goal: %w", err)
	}

	if goal == nil {
		return nil, errors.New("goal not found")
	}

	// Apply updates (Builder pattern for updates)
	updatedGoal := uc.applyUpdates(goal, request)

	// Validate updated goal
	if err := updatedGoal.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Save updates
	if err := uc.repository.Update(ctx, updatedGoal); err != nil {
		return nil, fmt.Errorf("failed to update goal: %w", err)
	}

	response := uc.convertGoalToUpdateResponse(updatedGoal)
	return response, nil
}

// DeleteGoal removes a savings goal
func (uc *SavingsGoalUseCase) DeleteGoal(ctx context.Context, request ports.DeleteSavingsGoalRequest) error {
	// Check if goal exists
	goal, err := uc.repository.GetByID(ctx, request.UserID, request.GoalID)
	if err != nil {
		return fmt.Errorf("failed to get goal: %w", err)
	}

	if goal == nil {
		return errors.New("goal not found")
	}

	// Business rule: Can't delete achieved goals with transactions
	if goal.Status == domain.SavingsGoalStatusAchieved && goal.CurrentAmount > 0 {
		return errors.New("cannot delete achieved goal with savings")
	}

	if err := uc.repository.Delete(ctx, request.UserID, request.GoalID); err != nil {
		return fmt.Errorf("failed to delete goal: %w", err)
	}

	return nil
}

// AddSavings adds money to a savings goal (Command pattern)
func (uc *SavingsGoalUseCase) AddSavings(ctx context.Context, request ports.AddSavingsRequest) (*ports.AddSavingsResponse, error) {
	// Get goal
	goal, err := uc.repository.GetByID(ctx, request.UserID, request.GoalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goal: %w", err)
	}

	if goal == nil {
		return nil, errors.New("goal not found")
	}

	// Business rules validation
	if goal.Status == domain.SavingsGoalStatusCancelled {
		return nil, errors.New("cannot add savings to cancelled goal")
	}

	if goal.Status == domain.SavingsGoalStatusPaused {
		return nil, errors.New("cannot add savings to paused goal")
	}

	// Create transaction
	transaction := &domain.SavingsTransaction{
		ID:          uuid.New().String(),
		UserID:      request.UserID,
		GoalID:      request.GoalID,
		Amount:      request.Amount,
		Type:        domain.SavingsTransactionTypeDeposit,
		Description: request.Description,
		CreatedAt:   time.Now(),
	}

	// Add savings to goal
	previousAmount := goal.CurrentAmount
	goal.AddSavings(request.Amount)

	// Check for milestone achievements
	uc.checkMilestones(ctx, goal, previousAmount)

	// Save transaction and update goal
	if err := uc.repository.CreateTransaction(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	if err := uc.repository.Update(ctx, goal); err != nil {
		return nil, fmt.Errorf("failed to update goal: %w", err)
	}

	// Send notifications if goal is achieved
	if goal.Status == domain.SavingsGoalStatusAchieved {
		if uc.notificationSvc != nil {
			if err := uc.notificationSvc.NotifyGoalAchieved(ctx, goal); err != nil {
				// Log error but don't fail the operation
				fmt.Printf("Failed to send achievement notification: %v\n", err)
			}
		}
	}

	response := &ports.AddSavingsResponse{
		TransactionID:    transaction.ID,
		GoalID:           goal.ID,
		Amount:           request.Amount,
		NewCurrentAmount: goal.CurrentAmount,
		NewProgress:      goal.Progress,
		RemainingAmount:  goal.RemainingAmount,
		IsAchieved:       goal.Status == domain.SavingsGoalStatusAchieved,
		Status:           goal.Status,
		CreatedAt:        transaction.CreatedAt,
	}

	return response, nil
}

// WithdrawSavings removes money from a savings goal
func (uc *SavingsGoalUseCase) WithdrawSavings(ctx context.Context, request ports.WithdrawSavingsRequest) (*ports.WithdrawSavingsResponse, error) {
	// Get goal
	goal, err := uc.repository.GetByID(ctx, request.UserID, request.GoalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goal: %w", err)
	}

	if goal == nil {
		return nil, errors.New("goal not found")
	}

	// Business rules validation
	if goal.Status == domain.SavingsGoalStatusCancelled {
		return nil, errors.New("cannot withdraw from cancelled goal")
	}

	if goal.CurrentAmount < request.Amount {
		return nil, errors.New("insufficient savings balance")
	}

	// Create transaction
	transaction := &domain.SavingsTransaction{
		ID:          uuid.New().String(),
		UserID:      request.UserID,
		GoalID:      request.GoalID,
		Amount:      request.Amount,
		Type:        domain.SavingsTransactionTypeWithdrawal,
		Description: request.Description,
		CreatedAt:   time.Now(),
	}

	// Withdraw from goal
	goal.WithdrawSavings(request.Amount)

	// Save transaction and update goal
	if err := uc.repository.CreateTransaction(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	if err := uc.repository.Update(ctx, goal); err != nil {
		return nil, fmt.Errorf("failed to update goal: %w", err)
	}

	response := &ports.WithdrawSavingsResponse{
		TransactionID:    transaction.ID,
		GoalID:           goal.ID,
		Amount:           request.Amount,
		NewCurrentAmount: goal.CurrentAmount,
		NewProgress:      goal.Progress,
		RemainingAmount:  goal.RemainingAmount,
		Status:           goal.Status,
		CreatedAt:        transaction.CreatedAt,
	}

	return response, nil
}

// PauseGoal pauses a savings goal
func (uc *SavingsGoalUseCase) PauseGoal(ctx context.Context, request ports.PauseGoalRequest) error {
	goal, err := uc.repository.GetByID(ctx, request.UserID, request.GoalID)
	if err != nil {
		return fmt.Errorf("failed to get goal: %w", err)
	}

	if goal == nil {
		return errors.New("goal not found")
	}

	if goal.Status != domain.SavingsGoalStatusActive {
		return errors.New("can only pause active goals")
	}

	goal.Pause()

	if err := uc.repository.Update(ctx, goal); err != nil {
		return fmt.Errorf("failed to pause goal: %w", err)
	}

	return nil
}

// ResumeGoal resumes a paused savings goal
func (uc *SavingsGoalUseCase) ResumeGoal(ctx context.Context, request ports.ResumeGoalRequest) error {
	goal, err := uc.repository.GetByID(ctx, request.UserID, request.GoalID)
	if err != nil {
		return fmt.Errorf("failed to get goal: %w", err)
	}

	if goal == nil {
		return errors.New("goal not found")
	}

	if goal.Status != domain.SavingsGoalStatusPaused {
		return errors.New("can only resume paused goals")
	}

	goal.Resume()

	if err := uc.repository.Update(ctx, goal); err != nil {
		return fmt.Errorf("failed to resume goal: %w", err)
	}

	return nil
}

// CancelGoal cancels a savings goal
func (uc *SavingsGoalUseCase) CancelGoal(ctx context.Context, request ports.CancelGoalRequest) error {
	goal, err := uc.repository.GetByID(ctx, request.UserID, request.GoalID)
	if err != nil {
		return fmt.Errorf("failed to get goal: %w", err)
	}

	if goal == nil {
		return errors.New("goal not found")
	}

	if goal.Status == domain.SavingsGoalStatusAchieved {
		return errors.New("cannot cancel achieved goal")
	}

	goal.Cancel()

	if err := uc.repository.Update(ctx, goal); err != nil {
		return fmt.Errorf("failed to cancel goal: %w", err)
	}

	return nil
}

// GetGoalSummary gets summary statistics for all goals
func (uc *SavingsGoalUseCase) GetGoalSummary(ctx context.Context, request ports.GetGoalSummaryRequest) (*ports.GetGoalSummaryResponse, error) {
	var goals []*domain.SavingsGoal
	var err error

	if request.Category != "" {
		goals, err = uc.repository.ListByCategory(ctx, request.UserID, request.Category)
	} else {
		goals, err = uc.repository.List(ctx, request.UserID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get goals: %w", err)
	}

	summary := uc.calculateSummary(goals)

	response := &ports.GetGoalSummaryResponse{
		Summary:            summary,
		MonthlyTargetTotal: uc.calculateMonthlyTargetTotal(goals),
		WeeklyTargetTotal:  uc.calculateWeeklyTargetTotal(goals),
		DailyTargetTotal:   uc.calculateDailyTargetTotal(goals),
		NextMilestones:     uc.calculateNextMilestones(goals),
		OverdueGoals:       uc.getOverdueGoalIDs(goals),
		UpdatedAt:          time.Now(),
	}

	return response, nil
}

// GetGoalTransactions gets transaction history for a goal
func (uc *SavingsGoalUseCase) GetGoalTransactions(ctx context.Context, request ports.GetGoalTransactionsRequest) (*ports.GetGoalTransactionsResponse, error) {
	// Validate goal exists and belongs to user
	goal, err := uc.repository.GetByID(ctx, request.UserID, request.GoalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goal: %w", err)
	}

	if goal == nil {
		return nil, errors.New("goal not found")
	}

	transactions, err := uc.repository.GetTransactionsByGoal(ctx, request.UserID, request.GoalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	// Apply pagination
	start := request.Offset
	end := start + request.Limit
	if request.Limit == 0 {
		end = len(transactions)
	}
	if start > len(transactions) {
		start = len(transactions)
	}
	if end > len(transactions) {
		end = len(transactions)
	}

	paginatedTransactions := transactions[start:end]

	response := &ports.GetGoalTransactionsResponse{
		Transactions: make([]domain.SavingsTransaction, len(paginatedTransactions)),
		Total:        len(transactions),
		Count:        len(paginatedTransactions),
	}

	// Convert pointers to values
	for i, txn := range paginatedTransactions {
		response.Transactions[i] = *txn
	}

	return response, nil
}

// Private helper methods

// validateCreateGoalRequest validates the create goal request
func (uc *SavingsGoalUseCase) validateCreateGoalRequest(request ports.CreateSavingsGoalRequest) error {
	if request.UserID == "" {
		return errors.New("user ID is required")
	}

	if request.Name == "" {
		return errors.New("goal name is required")
	}

	if request.TargetAmount <= 0 {
		return errors.New("target amount must be greater than 0")
	}

	if request.Category == "" {
		return errors.New("category is required")
	}

	// Target date (si viene) debe ser futura al crear
	if !request.TargetDate.IsZero() && request.TargetDate.Before(time.Now()) {
		return errors.New("target date must be in the future")
	}

	// Validate auto-save configuration
	if request.IsAutoSave {
		if request.AutoSaveAmount <= 0 {
			return errors.New("auto-save amount must be greater than 0")
		}
		if request.AutoSaveFrequency == "" {
			return errors.New("auto-save frequency is required when auto-save is enabled")
		}
	}

	return nil
}

// applyUpdates applies updates to a goal
func (uc *SavingsGoalUseCase) applyUpdates(goal *domain.SavingsGoal, request ports.UpdateSavingsGoalRequest) *domain.SavingsGoal {
	if request.Name != nil {
		goal.Name = *request.Name
	}
	if request.Description != nil {
		goal.Description = *request.Description
	}
	if request.TargetAmount != nil {
		goal.TargetAmount = *request.TargetAmount
		goal.UpdateCalculatedFields()
	}
	if request.Category != nil {
		goal.Category = *request.Category
	}
	if request.Priority != nil {
		goal.Priority = *request.Priority
	}
	if request.TargetDate != nil {
		goal.TargetDate = *request.TargetDate
		goal.UpdateCalculatedFields()
	}
	if request.IsAutoSave != nil {
		goal.IsAutoSave = *request.IsAutoSave
	}
	if request.AutoSaveAmount != nil {
		goal.AutoSaveAmount = *request.AutoSaveAmount
	}
	if request.AutoSaveFrequency != nil {
		goal.AutoSaveFrequency = *request.AutoSaveFrequency
	}
	if request.ImageURL != nil {
		goal.ImageURL = *request.ImageURL
	}

	goal.UpdatedAt = time.Now()
	return goal
}

// checkMilestones checks if any milestones were reached
func (uc *SavingsGoalUseCase) checkMilestones(ctx context.Context, goal *domain.SavingsGoal, previousAmount float64) {
	milestones := []float64{0.25, 0.50, 0.75}

	for _, milestone := range milestones {
		previousProgress := previousAmount / goal.TargetAmount
		currentProgress := goal.CurrentAmount / goal.TargetAmount

		if previousProgress < milestone && currentProgress >= milestone {
			// Milestone reached
			if uc.notificationSvc != nil {
				if err := uc.notificationSvc.NotifyMilestoneReached(ctx, goal, milestone); err != nil {
					fmt.Printf("Failed to send milestone notification: %v\n", err)
				}
			}
		}
	}
}

// calculateSummary calculates summary statistics for goals
func (uc *SavingsGoalUseCase) calculateSummary(goals []*domain.SavingsGoal) domain.SavingsGoalSummary {
	summary := domain.SavingsGoalSummary{}

	for _, goal := range goals {
		summary.TotalGoals++
		summary.TotalTarget += goal.TargetAmount
		summary.TotalSaved += goal.CurrentAmount

		switch goal.Status {
		case domain.SavingsGoalStatusActive:
			summary.ActiveGoals++
		case domain.SavingsGoalStatusAchieved:
			summary.AchievedGoals++
		case domain.SavingsGoalStatusPaused:
			summary.PausedGoals++
		case domain.SavingsGoalStatusCancelled:
			summary.CancelledGoals++
		}

		if goal.IsOverdue() {
			summary.OverdueGoals++
		}
	}

	if summary.TotalTarget > 0 {
		summary.AverageProgress = summary.TotalSaved / summary.TotalTarget
	}

	return summary
}

// calculateMonthlyTargetTotal calculates total monthly target for all active goals
func (uc *SavingsGoalUseCase) calculateMonthlyTargetTotal(goals []*domain.SavingsGoal) float64 {
	total := 0.0
	for _, goal := range goals {
		if goal.Status == domain.SavingsGoalStatusActive {
			total += goal.MonthlyTarget
		}
	}
	return total
}

// calculateWeeklyTargetTotal calculates total weekly target for all active goals
func (uc *SavingsGoalUseCase) calculateWeeklyTargetTotal(goals []*domain.SavingsGoal) float64 {
	total := 0.0
	for _, goal := range goals {
		if goal.Status == domain.SavingsGoalStatusActive {
			total += goal.WeeklyTarget
		}
	}
	return total
}

// calculateDailyTargetTotal calculates total daily target for all active goals
func (uc *SavingsGoalUseCase) calculateDailyTargetTotal(goals []*domain.SavingsGoal) float64 {
	total := 0.0
	for _, goal := range goals {
		if goal.Status == domain.SavingsGoalStatusActive {
			total += goal.DailyTarget
		}
	}
	return total
}

// calculateNextMilestones calculates next milestones for active goals
func (uc *SavingsGoalUseCase) calculateNextMilestones(goals []*domain.SavingsGoal) []ports.GoalMilestone {
	var milestones []ports.GoalMilestone

	for _, goal := range goals {
		if goal.Status != domain.SavingsGoalStatusActive {
			continue
		}

		// Find next milestone
		nextMilestone := uc.findNextMilestone(goal)
		if nextMilestone > 0 {
			amountNeeded := (goal.TargetAmount * nextMilestone) - goal.CurrentAmount
			estimatedDate := uc.estimateMilestoneDate(goal, amountNeeded)

			milestones = append(milestones, ports.GoalMilestone{
				GoalID:        goal.ID,
				GoalName:      goal.Name,
				Milestone:     nextMilestone,
				AmountNeeded:  amountNeeded,
				EstimatedDate: estimatedDate,
			})
		}
	}

	return milestones
}

// findNextMilestone finds the next milestone for a goal
func (uc *SavingsGoalUseCase) findNextMilestone(goal *domain.SavingsGoal) float64 {
	milestones := []float64{0.25, 0.50, 0.75, 1.0}

	for _, milestone := range milestones {
		if goal.Progress < milestone {
			return milestone
		}
	}

	return 0 // Already achieved all milestones
}

// estimateMilestoneDate estimates when a milestone will be reached
func (uc *SavingsGoalUseCase) estimateMilestoneDate(goal *domain.SavingsGoal, amountNeeded float64) time.Time {
	if goal.DailyTarget <= 0 {
		return time.Now().AddDate(0, 0, 365) // Default to 1 year if no target
	}

	daysNeeded := int(amountNeeded / goal.DailyTarget)
	return time.Now().AddDate(0, 0, daysNeeded)
}

// getOverdueGoalIDs gets IDs of overdue goals
func (uc *SavingsGoalUseCase) getOverdueGoalIDs(goals []*domain.SavingsGoal) []string {
	var overdueIDs []string

	for _, goal := range goals {
		if goal.IsOverdue() {
			overdueIDs = append(overdueIDs, goal.ID)
		}
	}

	return overdueIDs
}

// Conversion methods

// convertGoalToCreateResponse converts domain goal to create response
func (uc *SavingsGoalUseCase) convertGoalToCreateResponse(goal *domain.SavingsGoal) *ports.CreateSavingsGoalResponse {
	return &ports.CreateSavingsGoalResponse{
		ID:                goal.ID,
		UserID:            goal.UserID,
		Name:              goal.Name,
		Description:       goal.Description,
		TargetAmount:      goal.TargetAmount,
		CurrentAmount:     goal.CurrentAmount,
		Category:          goal.Category,
		Priority:          goal.Priority,
		TargetDate:        goal.TargetDate,
		Status:            goal.Status,
		MonthlyTarget:     goal.MonthlyTarget,
		WeeklyTarget:      goal.WeeklyTarget,
		DailyTarget:       goal.DailyTarget,
		Progress:          goal.Progress,
		DaysRemaining:     goal.DaysRemaining,
		IsAutoSave:        goal.IsAutoSave,
		AutoSaveAmount:    goal.AutoSaveAmount,
		AutoSaveFrequency: goal.AutoSaveFrequency,
		ImageURL:          goal.ImageURL,
		CreatedAt:         goal.CreatedAt,
		UpdatedAt:         goal.UpdatedAt,
	}
}

// convertGoalToGetResponse converts domain goal to get response
func (uc *SavingsGoalUseCase) convertGoalToGetResponse(goal *domain.SavingsGoal) *ports.GetSavingsGoalResponse {
	return &ports.GetSavingsGoalResponse{
		ID:                goal.ID,
		UserID:            goal.UserID,
		Name:              goal.Name,
		Description:       goal.Description,
		TargetAmount:      goal.TargetAmount,
		CurrentAmount:     goal.CurrentAmount,
		RemainingAmount:   goal.RemainingAmount,
		Category:          goal.Category,
		Priority:          goal.Priority,
		TargetDate:        goal.TargetDate,
		Status:            goal.Status,
		Progress:          goal.Progress,
		MonthlyTarget:     goal.MonthlyTarget,
		WeeklyTarget:      goal.WeeklyTarget,
		DailyTarget:       goal.DailyTarget,
		DaysRemaining:     goal.DaysRemaining,
		IsOverdue:         goal.IsOverdue(),
		IsOnTrack:         goal.IsOnTrack(),
		IsAutoSave:        goal.IsAutoSave,
		AutoSaveAmount:    goal.AutoSaveAmount,
		AutoSaveFrequency: goal.AutoSaveFrequency,
		ImageURL:          goal.ImageURL,
		CreatedAt:         goal.CreatedAt,
		UpdatedAt:         goal.UpdatedAt,
		AchievedAt:        goal.AchievedAt,
	}
}

// convertGoalToUpdateResponse converts domain goal to update response
func (uc *SavingsGoalUseCase) convertGoalToUpdateResponse(goal *domain.SavingsGoal) *ports.UpdateSavingsGoalResponse {
	return &ports.UpdateSavingsGoalResponse{
		ID:                goal.ID,
		UserID:            goal.UserID,
		Name:              goal.Name,
		Description:       goal.Description,
		TargetAmount:      goal.TargetAmount,
		CurrentAmount:     goal.CurrentAmount,
		RemainingAmount:   goal.RemainingAmount,
		Progress:          goal.Progress,
		Category:          goal.Category,
		Priority:          goal.Priority,
		TargetDate:        goal.TargetDate,
		Status:            goal.Status,
		MonthlyTarget:     goal.MonthlyTarget,
		WeeklyTarget:      goal.WeeklyTarget,
		DailyTarget:       goal.DailyTarget,
		IsAutoSave:        goal.IsAutoSave,
		AutoSaveAmount:    goal.AutoSaveAmount,
		AutoSaveFrequency: goal.AutoSaveFrequency,
		UpdatedAt:         goal.UpdatedAt,
	}
}
