package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/savings/domain"
)

// SavingsGoalModel is the GORM model for the savings_goals table
type SavingsGoalModel struct {
	ID                string     `gorm:"column:id;type:varchar(255);primaryKey"`
	UserID            string     `gorm:"column:user_id;type:varchar(255);not null;index"`
	TenantID          string     `gorm:"column:tenant_id;type:varchar(50);index"`
	Name              string     `gorm:"column:name;type:varchar(255);not null"`
	Description       string     `gorm:"column:description;type:text"`
	TargetAmount      float64    `gorm:"column:target_amount;not null"`
	CurrentAmount     float64    `gorm:"column:current_amount;not null;default:0"`
	Category          string     `gorm:"column:category;type:varchar(100);not null"`
	Priority          string     `gorm:"column:priority;type:varchar(50);not null"`
	TargetDate        time.Time  `gorm:"column:target_date;not null"`
	Status            string     `gorm:"column:status;type:varchar(50);not null;index"`
	MonthlyTarget     float64    `gorm:"column:monthly_target"`
	WeeklyTarget      float64    `gorm:"column:weekly_target"`
	DailyTarget       float64    `gorm:"column:daily_target"`
	Progress          float64    `gorm:"column:progress"`
	RemainingAmount   float64    `gorm:"column:remaining_amount"`
	DaysRemaining     int        `gorm:"column:days_remaining"`
	IsAutoSave        bool       `gorm:"column:is_auto_save;not null;default:false"`
	AutoSaveAmount    float64    `gorm:"column:auto_save_amount"`
	AutoSaveFrequency string     `gorm:"column:auto_save_frequency;type:varchar(50)"`
	ImageURL          string     `gorm:"column:image_url;type:text"`
	CreatedAt         time.Time  `gorm:"column:created_at"`
	UpdatedAt         time.Time  `gorm:"column:updated_at"`
	AchievedAt        *time.Time `gorm:"column:achieved_at"`
	DeletedAt         *time.Time `gorm:"column:deleted_at;index"`
}

func (SavingsGoalModel) TableName() string {
	return "savings_goals"
}

// ToSavingsGoal converts the GORM model to the domain entity
func (m *SavingsGoalModel) ToSavingsGoal() *domain.SavingsGoal {
	goal := &domain.SavingsGoal{
		ID:                m.ID,
		UserID:            m.UserID,
		TenantID:          m.TenantID,
		Name:              m.Name,
		Description:       m.Description,
		TargetAmount:      m.TargetAmount,
		CurrentAmount:     m.CurrentAmount,
		Category:          domain.SavingsGoalCategory(m.Category),
		Priority:          domain.SavingsGoalPriority(m.Priority),
		TargetDate:        m.TargetDate,
		Status:            domain.SavingsGoalStatus(m.Status),
		MonthlyTarget:     m.MonthlyTarget,
		WeeklyTarget:      m.WeeklyTarget,
		DailyTarget:       m.DailyTarget,
		Progress:          m.Progress,
		RemainingAmount:   m.RemainingAmount,
		DaysRemaining:     m.DaysRemaining,
		IsAutoSave:        m.IsAutoSave,
		AutoSaveAmount:    m.AutoSaveAmount,
		AutoSaveFrequency: m.AutoSaveFrequency,
		ImageURL:          m.ImageURL,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
		AchievedAt:        m.AchievedAt,
		DeletedAt:         m.DeletedAt,
	}
	return goal
}

// fromSavingsGoal converts the domain entity to the GORM model
func fromSavingsGoal(g *domain.SavingsGoal) *SavingsGoalModel {
	return &SavingsGoalModel{
		ID:                g.ID,
		UserID:            g.UserID,
		TenantID:          g.TenantID,
		Name:              g.Name,
		Description:       g.Description,
		TargetAmount:      g.TargetAmount,
		CurrentAmount:     g.CurrentAmount,
		Category:          string(g.Category),
		Priority:          string(g.Priority),
		TargetDate:        g.TargetDate,
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
		CreatedAt:         g.CreatedAt,
		UpdatedAt:         g.UpdatedAt,
		AchievedAt:        g.AchievedAt,
		DeletedAt:         g.DeletedAt,
	}
}

// SavingsTransactionModel is the GORM model for the savings_transactions table
type SavingsTransactionModel struct {
	ID          string    `gorm:"column:id;type:varchar(255);primaryKey"`
	GoalID      string    `gorm:"column:goal_id;type:varchar(255);not null;index"`
	UserID      string    `gorm:"column:user_id;type:varchar(255);not null;index"`
	Amount      float64   `gorm:"column:amount;not null"`
	Type        string    `gorm:"column:type;type:varchar(50);not null"`
	Description string    `gorm:"column:description;type:text"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

func (SavingsTransactionModel) TableName() string {
	return "savings_transactions"
}

// ToSavingsTransaction converts the GORM model to the domain entity
func (m *SavingsTransactionModel) ToSavingsTransaction() *domain.SavingsTransaction {
	return &domain.SavingsTransaction{
		ID:          m.ID,
		GoalID:      m.GoalID,
		UserID:      m.UserID,
		Amount:      m.Amount,
		Type:        domain.SavingsTransactionType(m.Type),
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
	}
}

// fromSavingsTransaction converts the domain entity to the GORM model
func fromSavingsTransaction(t *domain.SavingsTransaction) *SavingsTransactionModel {
	return &SavingsTransactionModel{
		ID:          t.ID,
		GoalID:      t.GoalID,
		UserID:      t.UserID,
		Amount:      t.Amount,
		Type:        string(t.Type),
		Description: t.Description,
		CreatedAt:   t.CreatedAt,
	}
}

// SavingsRepo implements ports.SavingsGoalRepository
type SavingsRepo struct {
	db *gorm.DB
}

// NewSavingsRepository creates a new SavingsRepo
func NewSavingsRepository(db *gorm.DB) *SavingsRepo {
	return &SavingsRepo{db: db}
}

// Create inserts a new savings goal into the database
func (r *SavingsRepo) Create(ctx context.Context, goal *domain.SavingsGoal) error {
	if goal.ID == "" {
		goal.ID = domain.NewSavingsGoalID()
	}
	if goal.ID == "" {
		goal.ID = "goal_" + uuid.New().String()[:8]
	}
	model := fromSavingsGoal(goal)
	return r.db.WithContext(ctx).Create(model).Error
}

// GetByID retrieves a savings goal by its ID, scoped to a tenant
func (r *SavingsRepo) GetByID(ctx context.Context, tenantID, goalID string) (*domain.SavingsGoal, error) {
	var model SavingsGoalModel
	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", goalID, tenantID).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.ToSavingsGoal(), nil
}

// List retrieves all savings goals for a tenant
func (r *SavingsRepo) List(ctx context.Context, tenantID string) ([]*domain.SavingsGoal, error) {
	var models []SavingsGoalModel
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID).
		Order("created_at DESC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	goals := make([]*domain.SavingsGoal, len(models))
	for i, m := range models {
		copy := m
		goals[i] = copy.ToSavingsGoal()
	}
	return goals, nil
}

// ListByStatus retrieves savings goals for a tenant filtered by status
func (r *SavingsRepo) ListByStatus(ctx context.Context, tenantID string, status domain.SavingsGoalStatus) ([]*domain.SavingsGoal, error) {
	var models []SavingsGoalModel
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND status = ? AND deleted_at IS NULL", tenantID, string(status)).
		Order("created_at DESC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	goals := make([]*domain.SavingsGoal, len(models))
	for i, m := range models {
		copy := m
		goals[i] = copy.ToSavingsGoal()
	}
	return goals, nil
}

// Update saves updated fields of an existing savings goal
func (r *SavingsRepo) Update(ctx context.Context, goal *domain.SavingsGoal) error {
	model := fromSavingsGoal(goal)
	return r.db.WithContext(ctx).
		Model(&SavingsGoalModel{}).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", goal.ID, goal.TenantID).
		Updates(model).Error
}

// Delete soft-deletes a savings goal
func (r *SavingsRepo) Delete(ctx context.Context, tenantID, goalID string) error {
	return r.db.WithContext(ctx).
		Model(&SavingsGoalModel{}).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", goalID, tenantID).
		Update("deleted_at", time.Now().UTC()).Error
}

// CreateTransaction inserts a savings transaction record
func (r *SavingsRepo) CreateTransaction(ctx context.Context, tx *domain.SavingsTransaction) error {
	if tx.ID == "" {
		tx.ID = domain.NewSavingsTransactionID()
	}
	model := fromSavingsTransaction(tx)
	return r.db.WithContext(ctx).Create(model).Error
}

// ListTransactions retrieves all transactions for a given savings goal
func (r *SavingsRepo) ListTransactions(ctx context.Context, goalID string) ([]*domain.SavingsTransaction, error) {
	var models []SavingsTransactionModel
	err := r.db.WithContext(ctx).
		Where("goal_id = ?", goalID).
		Order("created_at DESC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	txs := make([]*domain.SavingsTransaction, len(models))
	for i, m := range models {
		copy := m
		txs[i] = copy.ToSavingsTransaction()
	}
	return txs, nil
}
