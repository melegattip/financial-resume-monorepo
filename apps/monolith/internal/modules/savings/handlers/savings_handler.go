package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/savings/domain"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/savings/ports"
)

// SavingsHandler handles HTTP requests for savings goals
type SavingsHandler struct {
	repo   ports.SavingsGoalRepository
	logger zerolog.Logger
}

// NewSavingsHandler creates a new SavingsHandler
func NewSavingsHandler(repo ports.SavingsGoalRepository, logger zerolog.Logger) *SavingsHandler {
	return &SavingsHandler{
		repo:   repo,
		logger: logger,
	}
}

// --- Request / Response types ---

// CreateGoalRequest is the request body for creating a savings goal
type CreateGoalRequest struct {
	Name              string  `json:"name" binding:"required"`
	Description       string  `json:"description"`
	TargetAmount      float64 `json:"target_amount" binding:"required,gt=0"`
	Category          string  `json:"category" binding:"required"`
	Priority          string  `json:"priority"`
	TargetDate        string  `json:"target_date" binding:"required"` // RFC3339
	IsAutoSave        bool    `json:"is_auto_save"`
	AutoSaveAmount    float64 `json:"auto_save_amount"`
	AutoSaveFrequency string  `json:"auto_save_frequency"`
	ImageURL          string  `json:"image_url"`
}

// UpdateGoalRequest is the request body for updating a savings goal
type UpdateGoalRequest struct {
	Name              string  `json:"name" binding:"required"`
	Description       string  `json:"description"`
	TargetAmount      float64 `json:"target_amount" binding:"required,gt=0"`
	Category          string  `json:"category" binding:"required"`
	Priority          string  `json:"priority"`
	TargetDate        string  `json:"target_date" binding:"required"` // RFC3339
	IsAutoSave        bool    `json:"is_auto_save"`
	AutoSaveAmount    float64 `json:"auto_save_amount"`
	AutoSaveFrequency string  `json:"auto_save_frequency"`
	ImageURL          string  `json:"image_url"`
}

// AmountRequest is the body for deposit/withdraw operations
type AmountRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description"`
}

// SavingsGoalResponse is the HTTP response for a savings goal
type SavingsGoalResponse struct {
	ID                string  `json:"id"`
	UserID            string  `json:"user_id"`
	Name              string  `json:"name"`
	Description       string  `json:"description,omitempty"`
	TargetAmount      float64 `json:"target_amount"`
	CurrentAmount     float64 `json:"current_amount"`
	Category          string  `json:"category"`
	Priority          string  `json:"priority"`
	TargetDate        string  `json:"target_date"`
	Status            string  `json:"status"`
	MonthlyTarget     float64 `json:"monthly_target"`
	WeeklyTarget      float64 `json:"weekly_target"`
	DailyTarget       float64 `json:"daily_target"`
	Progress          float64 `json:"progress"`
	RemainingAmount   float64 `json:"remaining_amount"`
	DaysRemaining     int     `json:"days_remaining"`
	IsAutoSave        bool    `json:"is_auto_save"`
	AutoSaveAmount    float64 `json:"auto_save_amount,omitempty"`
	AutoSaveFrequency string  `json:"auto_save_frequency,omitempty"`
	ImageURL          string  `json:"image_url,omitempty"`
	IsOverdue         bool    `json:"is_overdue"`
	IsOnTrack         bool    `json:"is_on_track"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
	AchievedAt        *string `json:"achieved_at,omitempty"`
}

// SavingsTransactionResponse is the HTTP response for a savings transaction
type SavingsTransactionResponse struct {
	ID          string  `json:"id"`
	GoalID      string  `json:"goal_id"`
	UserID      string  `json:"user_id"`
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"`
	Description string  `json:"description,omitempty"`
	CreatedAt   string  `json:"created_at"`
}

func toGoalResponse(g *domain.SavingsGoal) SavingsGoalResponse {
	resp := SavingsGoalResponse{
		ID:                g.ID,
		UserID:            g.UserID,
		Name:              g.Name,
		Description:       g.Description,
		TargetAmount:      g.TargetAmount,
		CurrentAmount:     g.CurrentAmount,
		Category:          string(g.Category),
		Priority:          string(g.Priority),
		TargetDate:        g.TargetDate.Format(time.RFC3339),
		Status:            string(g.Status),
		MonthlyTarget:     g.MonthlyTarget,
		WeeklyTarget:      g.WeeklyTarget,
		DailyTarget:       g.DailyTarget,
		Progress:          g.Progress,
		RemainingAmount:   g.RemainingAmount,
		DaysRemaining:     g.DaysRemaining,
		IsAutoSave:        g.IsAutoSave,
		AutoSaveAmount:    g.AutoSaveAmount,
		AutoSaveFrequency: g.AutoSaveFrequency,
		ImageURL:          g.ImageURL,
		IsOverdue:         g.IsOverdue(),
		IsOnTrack:         g.IsOnTrack(),
		CreatedAt:         g.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         g.UpdatedAt.Format(time.RFC3339),
	}
	if g.AchievedAt != nil {
		s := g.AchievedAt.Format(time.RFC3339)
		resp.AchievedAt = &s
	}
	return resp
}

func toTransactionResponse(t *domain.SavingsTransaction) SavingsTransactionResponse {
	return SavingsTransactionResponse{
		ID:          t.ID,
		GoalID:      t.GoalID,
		UserID:      t.UserID,
		Amount:      t.Amount,
		Type:        string(t.Type),
		Description: t.Description,
		CreatedAt:   t.CreatedAt.Format(time.RFC3339),
	}
}

// --- Helpers ---

func getUserID(c *gin.Context) (string, bool) {
	val, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	id, ok := val.(string)
	return id, ok
}

// --- Handlers ---

// CreateGoal handles POST /savings
func (h *SavingsHandler) CreateGoal(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	targetDate, err := time.Parse(time.RFC3339, req.TargetDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target_date format, expected RFC3339"})
		return
	}

	if !targetDate.After(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrInvalidTargetDate.Error()})
		return
	}

	priority := domain.SavingsGoalPriorityMedium
	if req.Priority != "" {
		priority = domain.SavingsGoalPriority(req.Priority)
	}

	builder := domain.NewSavingsGoalBuilder().
		SetUserID(userID).
		SetName(req.Name).
		SetDescription(req.Description).
		SetTargetAmount(req.TargetAmount).
		SetCategory(domain.SavingsGoalCategory(req.Category)).
		SetPriority(priority).
		SetTargetDate(targetDate).
		SetImageURL(req.ImageURL)

	if req.IsAutoSave {
		builder.SetAutoSave(req.AutoSaveAmount, req.AutoSaveFrequency)
	}

	goal, err := builder.Build()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Create(c.Request.Context(), goal); err != nil {
		h.logger.Error().Err(err).Msg("failed to create savings goal")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create savings goal"})
		return
	}

	goal.UpdateCalculatedFields()
	c.JSON(http.StatusCreated, toGoalResponse(goal))
}

// ListGoals handles GET /savings
func (h *SavingsHandler) ListGoals(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	statusFilter := c.Query("status")

	var (
		goals []*domain.SavingsGoal
		err   error
	)

	if statusFilter != "" {
		goals, err = h.repo.ListByStatus(c.Request.Context(), userID, domain.SavingsGoalStatus(statusFilter))
	} else {
		goals, err = h.repo.List(c.Request.Context(), userID)
	}

	if err != nil {
		h.logger.Error().Err(err).Msg("failed to list savings goals")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list savings goals"})
		return
	}

	// Apply optional client-side filters (category, priority) and update calculated fields
	categoryFilter := c.Query("category")
	priorityFilter := c.Query("priority")

	response := make([]SavingsGoalResponse, 0, len(goals))
	for _, g := range goals {
		if categoryFilter != "" && string(g.Category) != categoryFilter {
			continue
		}
		if priorityFilter != "" && string(g.Priority) != priorityFilter {
			continue
		}
		g.UpdateCalculatedFields()
		response = append(response, toGoalResponse(g))
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  response,
		"total": len(response),
	})
}

// GetSummary handles GET /savings/summary
func (h *SavingsHandler) GetSummary(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	goals, err := h.repo.List(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to fetch savings goals for summary")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch savings summary"})
		return
	}

	summary := &domain.SavingsGoalSummary{}
	totalProgress := 0.0

	for _, g := range goals {
		g.UpdateCalculatedFields()
		summary.TotalGoals++
		summary.TotalTarget += g.TargetAmount
		summary.TotalSaved += g.CurrentAmount
		summary.TotalRemaining += g.GetRemainingAmount()
		totalProgress += g.GetProgress()

		switch g.Status {
		case domain.SavingsGoalStatusActive:
			summary.ActiveGoals++
		case domain.SavingsGoalStatusAchieved:
			summary.AchievedGoals++
		case domain.SavingsGoalStatusPaused:
			summary.PausedGoals++
		case domain.SavingsGoalStatusCancelled:
			summary.CancelledGoals++
		}

		if g.IsOverdue() {
			summary.OverdueGoals++
		}
		if g.IsOnTrack() {
			summary.OnTrackGoals++
		}
	}

	if summary.TotalGoals > 0 {
		summary.AverageProgress = totalProgress / float64(summary.TotalGoals)
	}

	c.JSON(http.StatusOK, summary)
}

// GetGoal handles GET /savings/:id
func (h *SavingsHandler) GetGoal(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	goalID := c.Param("id")
	goal, err := h.repo.GetByID(c.Request.Context(), userID, goalID)
	if err != nil {
		h.logger.Error().Err(err).Str("goal_id", goalID).Msg("failed to get savings goal")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get savings goal"})
		return
	}
	if goal == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "savings goal not found"})
		return
	}

	goal.UpdateCalculatedFields()
	c.JSON(http.StatusOK, toGoalResponse(goal))
}

// UpdateGoal handles PUT /savings/:id
func (h *SavingsHandler) UpdateGoal(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	goalID := c.Param("id")
	goal, err := h.repo.GetByID(c.Request.Context(), userID, goalID)
	if err != nil {
		h.logger.Error().Err(err).Str("goal_id", goalID).Msg("failed to get savings goal")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get savings goal"})
		return
	}
	if goal == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "savings goal not found"})
		return
	}

	var req UpdateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	targetDate, err := time.Parse(time.RFC3339, req.TargetDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target_date format, expected RFC3339"})
		return
	}

	// Apply updates
	goal.Name = req.Name
	goal.Description = req.Description
	goal.TargetAmount = req.TargetAmount
	goal.Category = domain.SavingsGoalCategory(req.Category)
	if req.Priority != "" {
		goal.Priority = domain.SavingsGoalPriority(req.Priority)
	}
	goal.TargetDate = targetDate
	goal.IsAutoSave = req.IsAutoSave
	goal.AutoSaveAmount = req.AutoSaveAmount
	goal.AutoSaveFrequency = req.AutoSaveFrequency
	goal.ImageURL = req.ImageURL
	goal.UpdatedAt = time.Now().UTC()

	if err := goal.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goal.UpdateCalculatedFields()

	if err := h.repo.Update(c.Request.Context(), goal); err != nil {
		h.logger.Error().Err(err).Msg("failed to update savings goal")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update savings goal"})
		return
	}

	c.JSON(http.StatusOK, toGoalResponse(goal))
}

// DeleteGoal handles DELETE /savings/:id
func (h *SavingsHandler) DeleteGoal(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	goalID := c.Param("id")
	goal, err := h.repo.GetByID(c.Request.Context(), userID, goalID)
	if err != nil {
		h.logger.Error().Err(err).Str("goal_id", goalID).Msg("failed to get savings goal")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get savings goal"})
		return
	}
	if goal == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "savings goal not found"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), userID, goalID); err != nil {
		h.logger.Error().Err(err).Msg("failed to delete savings goal")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete savings goal"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// AddSavings handles POST /savings/:id/deposit
func (h *SavingsHandler) AddSavings(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	goalID := c.Param("id")
	goal, err := h.repo.GetByID(c.Request.Context(), userID, goalID)
	if err != nil {
		h.logger.Error().Err(err).Str("goal_id", goalID).Msg("failed to get savings goal")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get savings goal"})
		return
	}
	if goal == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "savings goal not found"})
		return
	}

	var req AmountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := goal.AddSavings(req.Amount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Update(c.Request.Context(), goal); err != nil {
		h.logger.Error().Err(err).Msg("failed to update savings goal after deposit")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update savings goal"})
		return
	}

	tx := &domain.SavingsTransaction{
		GoalID:      goalID,
		UserID:      userID,
		Amount:      req.Amount,
		Type:        domain.SavingsTransactionTypeDeposit,
		Description: req.Description,
		CreatedAt:   time.Now().UTC(),
	}
	if err := h.repo.CreateTransaction(c.Request.Context(), tx); err != nil {
		h.logger.Warn().Err(err).Msg("failed to record deposit transaction")
	}

	goal.UpdateCalculatedFields()
	c.JSON(http.StatusOK, toGoalResponse(goal))
}

// WithdrawSavings handles POST /savings/:id/withdraw
func (h *SavingsHandler) WithdrawSavings(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	goalID := c.Param("id")
	goal, err := h.repo.GetByID(c.Request.Context(), userID, goalID)
	if err != nil {
		h.logger.Error().Err(err).Str("goal_id", goalID).Msg("failed to get savings goal")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get savings goal"})
		return
	}
	if goal == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "savings goal not found"})
		return
	}

	var req AmountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := goal.WithdrawSavings(req.Amount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Update(c.Request.Context(), goal); err != nil {
		h.logger.Error().Err(err).Msg("failed to update savings goal after withdrawal")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update savings goal"})
		return
	}

	tx := &domain.SavingsTransaction{
		GoalID:      goalID,
		UserID:      userID,
		Amount:      req.Amount,
		Type:        domain.SavingsTransactionTypeWithdrawal,
		Description: req.Description,
		CreatedAt:   time.Now().UTC(),
	}
	if err := h.repo.CreateTransaction(c.Request.Context(), tx); err != nil {
		h.logger.Warn().Err(err).Msg("failed to record withdrawal transaction")
	}

	goal.UpdateCalculatedFields()
	c.JSON(http.StatusOK, toGoalResponse(goal))
}

// PauseGoal handles POST /savings/:id/pause
func (h *SavingsHandler) PauseGoal(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	goalID := c.Param("id")
	goal, err := h.repo.GetByID(c.Request.Context(), userID, goalID)
	if err != nil {
		h.logger.Error().Err(err).Str("goal_id", goalID).Msg("failed to get savings goal")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get savings goal"})
		return
	}
	if goal == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "savings goal not found"})
		return
	}

	if err := goal.Pause(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Update(c.Request.Context(), goal); err != nil {
		h.logger.Error().Err(err).Msg("failed to pause savings goal")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to pause savings goal"})
		return
	}

	goal.UpdateCalculatedFields()
	c.JSON(http.StatusOK, toGoalResponse(goal))
}

// ResumeGoal handles POST /savings/:id/resume
func (h *SavingsHandler) ResumeGoal(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	goalID := c.Param("id")
	goal, err := h.repo.GetByID(c.Request.Context(), userID, goalID)
	if err != nil {
		h.logger.Error().Err(err).Str("goal_id", goalID).Msg("failed to get savings goal")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get savings goal"})
		return
	}
	if goal == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "savings goal not found"})
		return
	}

	if err := goal.Resume(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Update(c.Request.Context(), goal); err != nil {
		h.logger.Error().Err(err).Msg("failed to resume savings goal")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resume savings goal"})
		return
	}

	goal.UpdateCalculatedFields()
	c.JSON(http.StatusOK, toGoalResponse(goal))
}

// CancelGoal handles POST /savings/:id/cancel
func (h *SavingsHandler) CancelGoal(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	goalID := c.Param("id")
	goal, err := h.repo.GetByID(c.Request.Context(), userID, goalID)
	if err != nil {
		h.logger.Error().Err(err).Str("goal_id", goalID).Msg("failed to get savings goal")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get savings goal"})
		return
	}
	if goal == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "savings goal not found"})
		return
	}

	if err := goal.Cancel(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Update(c.Request.Context(), goal); err != nil {
		h.logger.Error().Err(err).Msg("failed to cancel savings goal")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cancel savings goal"})
		return
	}

	goal.UpdateCalculatedFields()
	c.JSON(http.StatusOK, toGoalResponse(goal))
}

// ListTransactions handles GET /savings/:id/transactions
func (h *SavingsHandler) ListTransactions(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	goalID := c.Param("id")

	// Verify ownership
	goal, err := h.repo.GetByID(c.Request.Context(), userID, goalID)
	if err != nil {
		h.logger.Error().Err(err).Str("goal_id", goalID).Msg("failed to get savings goal")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get savings goal"})
		return
	}
	if goal == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "savings goal not found"})
		return
	}

	txs, err := h.repo.ListTransactions(c.Request.Context(), goalID)
	if err != nil {
		h.logger.Error().Err(err).Str("goal_id", goalID).Msg("failed to list savings transactions")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list transactions"})
		return
	}

	response := make([]SavingsTransactionResponse, len(txs))
	for i, t := range txs {
		response[i] = toTransactionResponse(t)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  response,
		"total": len(response),
	})
}
