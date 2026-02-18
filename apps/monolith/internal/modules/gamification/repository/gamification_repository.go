package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/gamification/domain"
)

// ---------------------------------------------------------------------------
// GORM models
// ---------------------------------------------------------------------------

// UserGamificationModel is the GORM model for the user_gamification table.
type UserGamificationModel struct {
	ID                string     `gorm:"column:id;type:varchar(255);primaryKey"`
	UserID            string     `gorm:"column:user_id;type:varchar(255);not null;uniqueIndex"`
	TotalXP           int        `gorm:"column:total_xp;not null;default:0"`
	CurrentLevel      int        `gorm:"column:current_level;not null;default:1"`
	InsightsViewed    int        `gorm:"column:insights_viewed;not null;default:0"`
	ActionsCompleted  int        `gorm:"column:actions_completed;not null;default:0"`
	AchievementsCount int        `gorm:"column:achievements_count;not null;default:0"`
	CurrentStreak     int        `gorm:"column:current_streak;not null;default:0"`
	LastActivity      time.Time  `gorm:"column:last_activity"`
	CreatedAt         time.Time  `gorm:"column:created_at"`
	UpdatedAt         time.Time  `gorm:"column:updated_at"`
	DeletedAt         *time.Time `gorm:"column:deleted_at;index"`
}

func (UserGamificationModel) TableName() string {
	return "user_gamification"
}

func (m *UserGamificationModel) toDomain() *domain.UserGamification {
	return &domain.UserGamification{
		ID:                m.ID,
		UserID:            m.UserID,
		TotalXP:           m.TotalXP,
		CurrentLevel:      m.CurrentLevel,
		InsightsViewed:    m.InsightsViewed,
		ActionsCompleted:  m.ActionsCompleted,
		AchievementsCount: m.AchievementsCount,
		CurrentStreak:     m.CurrentStreak,
		LastActivity:      m.LastActivity,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
		DeletedAt:         m.DeletedAt,
	}
}

func fromDomainGamification(g *domain.UserGamification) *UserGamificationModel {
	return &UserGamificationModel{
		ID:                g.ID,
		UserID:            g.UserID,
		TotalXP:           g.TotalXP,
		CurrentLevel:      g.CurrentLevel,
		InsightsViewed:    g.InsightsViewed,
		ActionsCompleted:  g.ActionsCompleted,
		AchievementsCount: g.AchievementsCount,
		CurrentStreak:     g.CurrentStreak,
		LastActivity:      g.LastActivity,
		CreatedAt:         g.CreatedAt,
		UpdatedAt:         g.UpdatedAt,
		DeletedAt:         g.DeletedAt,
	}
}

// AchievementModel is the GORM model for the achievements table.
type AchievementModel struct {
	ID          string     `gorm:"column:id;type:varchar(255);primaryKey"`
	UserID      string     `gorm:"column:user_id;type:varchar(255);not null;index"`
	Type        string     `gorm:"column:type;type:varchar(100);not null"`
	Name        string     `gorm:"column:name;type:varchar(255);not null"`
	Description string     `gorm:"column:description;type:text"`
	Points      int        `gorm:"column:points;not null;default:0"`
	Progress    int        `gorm:"column:progress;not null;default:0"`
	Target      int        `gorm:"column:target;not null;default:1"`
	Completed   bool       `gorm:"column:completed;not null;default:false"`
	UnlockedAt  *time.Time `gorm:"column:unlocked_at"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
}

func (AchievementModel) TableName() string {
	return "achievements"
}

func (m *AchievementModel) toDomain() *domain.Achievement {
	return &domain.Achievement{
		ID:          m.ID,
		UserID:      m.UserID,
		Type:        m.Type,
		Name:        m.Name,
		Description: m.Description,
		Points:      m.Points,
		Progress:    m.Progress,
		Target:      m.Target,
		Completed:   m.Completed,
		UnlockedAt:  m.UnlockedAt,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func fromDomainAchievement(a *domain.Achievement) *AchievementModel {
	return &AchievementModel{
		ID:          a.ID,
		UserID:      a.UserID,
		Type:        a.Type,
		Name:        a.Name,
		Description: a.Description,
		Points:      a.Points,
		Progress:    a.Progress,
		Target:      a.Target,
		Completed:   a.Completed,
		UnlockedAt:  a.UnlockedAt,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

// UserActionModel is the GORM model for the user_actions table.
// This table intentionally has no deleted_at column.
type UserActionModel struct {
	ID          string    `gorm:"column:id;type:varchar(255);primaryKey"`
	UserID      string    `gorm:"column:user_id;type:varchar(255);not null;index"`
	ActionType  string    `gorm:"column:action_type;type:varchar(100);not null"`
	EntityType  string    `gorm:"column:entity_type;type:varchar(100)"`
	EntityID    string    `gorm:"column:entity_id;type:varchar(255)"`
	Description string    `gorm:"column:description;type:text"`
	XPEarned    int       `gorm:"column:xp_earned;not null;default:0"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

func (UserActionModel) TableName() string {
	return "user_actions"
}

func (m *UserActionModel) toDomain() *domain.UserAction {
	return &domain.UserAction{
		ID:          m.ID,
		UserID:      m.UserID,
		ActionType:  m.ActionType,
		EntityType:  m.EntityType,
		EntityID:    m.EntityID,
		Description: m.Description,
		XPEarned:    m.XPEarned,
		CreatedAt:   m.CreatedAt,
	}
}

func fromDomainUserAction(a *domain.UserAction) *UserActionModel {
	return &UserActionModel{
		ID:          a.ID,
		UserID:      a.UserID,
		ActionType:  a.ActionType,
		EntityType:  a.EntityType,
		EntityID:    a.EntityID,
		Description: a.Description,
		XPEarned:    a.XPEarned,
		CreatedAt:   a.CreatedAt,
	}
}

// ---------------------------------------------------------------------------
// Repository
// ---------------------------------------------------------------------------

// GamificationRepo provides persistence for the gamification module.
type GamificationRepo struct {
	db *gorm.DB
}

// NewGamificationRepository creates a new GamificationRepo.
func NewGamificationRepository(db *gorm.DB) *GamificationRepo {
	return &GamificationRepo{db: db}
}

// CreateUserGamification persists a new UserGamification aggregate.
func (r *GamificationRepo) CreateUserGamification(ctx context.Context, g *domain.UserGamification) error {
	if g.ID == "" {
		g.ID = uuid.New().String()
	}
	model := fromDomainGamification(g)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	g.ID = model.ID
	return nil
}

// FindUserGamificationByUserID retrieves the gamification record for a user.
// Returns nil, nil when no record exists.
func (r *GamificationRepo) FindUserGamificationByUserID(ctx context.Context, userID string) (*domain.UserGamification, error) {
	var model UserGamificationModel
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.toDomain(), nil
}

// UpdateUserGamification saves all changes to an existing UserGamification record.
func (r *GamificationRepo) UpdateUserGamification(ctx context.Context, g *domain.UserGamification) error {
	model := fromDomainGamification(g)
	return r.db.WithContext(ctx).
		Model(&UserGamificationModel{}).
		Where("id = ? AND deleted_at IS NULL", g.ID).
		Updates(model).Error
}

// CreateAchievements bulk-inserts a slice of Achievement records.
func (r *GamificationRepo) CreateAchievements(ctx context.Context, achievements []domain.Achievement) error {
	if len(achievements) == 0 {
		return nil
	}
	models := make([]AchievementModel, len(achievements))
	for i, a := range achievements {
		if a.ID == "" {
			a.ID = uuid.New().String()
		}
		models[i] = *fromDomainAchievement(&a)
	}
	return r.db.WithContext(ctx).Create(&models).Error
}

// FindAchievementsByUserID retrieves all achievements for a user.
func (r *GamificationRepo) FindAchievementsByUserID(ctx context.Context, userID string) ([]domain.Achievement, error) {
	var models []AchievementModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&models).Error; err != nil {
		return nil, err
	}
	achievements := make([]domain.Achievement, len(models))
	for i, m := range models {
		achievements[i] = *m.toDomain()
	}
	return achievements, nil
}

// UpdateAchievement saves changes to an existing Achievement record.
func (r *GamificationRepo) UpdateAchievement(ctx context.Context, a *domain.Achievement) error {
	model := fromDomainAchievement(a)
	return r.db.WithContext(ctx).
		Model(&AchievementModel{}).
		Where("id = ?", a.ID).
		Updates(model).Error
}

// CreateUserAction persists a new UserAction record.
func (r *GamificationRepo) CreateUserAction(ctx context.Context, a *domain.UserAction) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	model := fromDomainUserAction(a)
	return r.db.WithContext(ctx).Create(model).Error
}

// FindActionsByUserID retrieves all actions recorded for a user.
func (r *GamificationRepo) FindActionsByUserID(ctx context.Context, userID string) ([]domain.UserAction, error) {
	var models []UserActionModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}
	actions := make([]domain.UserAction, len(models))
	for i, m := range models {
		actions[i] = *m.toDomain()
	}
	return actions, nil
}

// FindActionsByUserIDAndDay retrieves actions for a user within the given calendar day.
// This is used for idempotency checks (e.g., only award view_dashboard XP once per day).
func (r *GamificationRepo) FindActionsByUserIDAndDay(ctx context.Context, userID string, day time.Time) ([]domain.UserAction, error) {
	startOfDay := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	var models []UserActionModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND created_at >= ? AND created_at < ?", userID, startOfDay, endOfDay).
		Find(&models).Error; err != nil {
		return nil, err
	}
	actions := make([]domain.UserAction, len(models))
	for i, m := range models {
		actions[i] = *m.toDomain()
	}
	return actions, nil
}
