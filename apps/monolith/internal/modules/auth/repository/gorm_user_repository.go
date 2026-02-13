package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/domain"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/ports"
)

// gormRepository implements all auth repository interfaces using GORM.
type gormRepository struct {
	db *gorm.DB
}

// New creates a new GORM-backed repository that implements all auth repository interfaces.
func New(db *gorm.DB) *gormRepository {
	return &gormRepository{db: db}
}

// Compile-time interface checks.
var (
	_ ports.UserRepository                = (*gormRepository)(nil)
	_ ports.PreferencesRepository         = (*gormRepository)(nil)
	_ ports.NotificationSettingsRepository = (*gormRepository)(nil)
	_ ports.TwoFARepository               = (*gormRepository)(nil)
)

// ---------------------------------------------------------------------------
// UserRepository
// ---------------------------------------------------------------------------

func (r *gormRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return &user, nil
}

func (r *gormRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	return &user, nil
}

func (r *gormRepository) FindByVerificationToken(ctx context.Context, token string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Where("email_verification_token = ? AND email_verification_expires > ?", token, time.Now()).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("invalid or expired email verification token")
		}
		return nil, fmt.Errorf("failed to find user by verification token: %w", err)
	}
	return &user, nil
}

func (r *gormRepository) FindByResetToken(ctx context.Context, token string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Where("password_reset_token = ? AND password_reset_expires > ?", token, time.Now()).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("invalid or expired password reset token")
		}
		return nil, fmt.Errorf("failed to find user by reset token: %w", err)
	}
	return &user, nil
}

func (r *gormRepository) Create(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if isDuplicateKeyError(err) {
			return fmt.Errorf("user with email %s already exists", user.Email)
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *gormRepository) Update(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *gormRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.User{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user with id %s not found", id)
	}
	return nil
}

// ---------------------------------------------------------------------------
// PreferencesRepository
// ---------------------------------------------------------------------------

func (r *gormRepository) FindPreferencesByUserID(ctx context.Context, userID string) (*domain.Preferences, error) {
	var prefs domain.Preferences
	if err := r.db.WithContext(ctx).First(&prefs, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("preferences for user %s not found", userID)
		}
		return nil, fmt.Errorf("failed to find preferences: %w", err)
	}
	return &prefs, nil
}

func (r *gormRepository) CreateDefaultPreferences(ctx context.Context, userID string) (*domain.Preferences, error) {
	prefs := domain.DefaultPreferences(userID)
	if err := r.db.WithContext(ctx).Create(&prefs).Error; err != nil {
		return nil, fmt.Errorf("failed to create default preferences: %w", err)
	}
	return &prefs, nil
}

// UpdatePreferences updates existing preferences.
func (r *gormRepository) UpdatePreferences(ctx context.Context, prefs *domain.Preferences) error {
	if err := r.db.WithContext(ctx).Save(prefs).Error; err != nil {
		return fmt.Errorf("failed to update preferences: %w", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// NotificationSettingsRepository
// ---------------------------------------------------------------------------

func (r *gormRepository) FindNotificationsByUserID(ctx context.Context, userID string) (*domain.NotificationSettings, error) {
	var settings domain.NotificationSettings
	if err := r.db.WithContext(ctx).First(&settings, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("notification settings for user %s not found", userID)
		}
		return nil, fmt.Errorf("failed to find notification settings: %w", err)
	}
	return &settings, nil
}

func (r *gormRepository) CreateDefaultNotifications(ctx context.Context, userID string) (*domain.NotificationSettings, error) {
	settings := domain.DefaultNotificationSettings(userID)
	if err := r.db.WithContext(ctx).Create(&settings).Error; err != nil {
		return nil, fmt.Errorf("failed to create default notification settings: %w", err)
	}
	return &settings, nil
}

func (r *gormRepository) UpdateNotifications(ctx context.Context, settings *domain.NotificationSettings) error {
	if err := r.db.WithContext(ctx).Save(settings).Error; err != nil {
		return fmt.Errorf("failed to update notification settings: %w", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// TwoFARepository
// ---------------------------------------------------------------------------

func (r *gormRepository) FindTwoFAByUserID(ctx context.Context, userID string) (*domain.TwoFA, error) {
	var twoFA domain.TwoFA
	if err := r.db.WithContext(ctx).First(&twoFA, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("2FA settings for user %s not found", userID)
		}
		return nil, fmt.Errorf("failed to find 2FA settings: %w", err)
	}
	return &twoFA, nil
}

func (r *gormRepository) UpsertTwoFA(ctx context.Context, twoFA *domain.TwoFA) error {
	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"secret", "enabled", "backup_codes", "last_used_code", "updated_at"}),
		}).
		Create(twoFA).Error
	if err != nil {
		return fmt.Errorf("failed to upsert 2FA settings: %w", err)
	}
	return nil
}

func (r *gormRepository) DeleteTwoFA(ctx context.Context, userID string) error {
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&domain.TwoFA{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete 2FA settings: %w", result.Error)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func isDuplicateKeyError(err error) bool {
	return strings.Contains(err.Error(), "23505") || strings.Contains(err.Error(), "duplicate key")
}
