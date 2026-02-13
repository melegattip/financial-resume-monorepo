package ports

import (
	"context"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/domain"
)

// UserRepository defines persistence operations for User entities.
type UserRepository interface {
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByVerificationToken(ctx context.Context, token string) (*domain.User, error)
	FindByResetToken(ctx context.Context, token string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
}

// PreferencesRepository defines persistence operations for Preferences.
type PreferencesRepository interface {
	FindPreferencesByUserID(ctx context.Context, userID string) (*domain.Preferences, error)
	CreateDefaultPreferences(ctx context.Context, userID string) (*domain.Preferences, error)
	UpdatePreferences(ctx context.Context, prefs *domain.Preferences) error
}

// NotificationSettingsRepository defines persistence operations for NotificationSettings.
type NotificationSettingsRepository interface {
	FindNotificationsByUserID(ctx context.Context, userID string) (*domain.NotificationSettings, error)
	CreateDefaultNotifications(ctx context.Context, userID string) (*domain.NotificationSettings, error)
	UpdateNotifications(ctx context.Context, settings *domain.NotificationSettings) error
}

// TwoFARepository defines persistence operations for TwoFA.
type TwoFARepository interface {
	FindTwoFAByUserID(ctx context.Context, userID string) (*domain.TwoFA, error)
	UpsertTwoFA(ctx context.Context, twoFA *domain.TwoFA) error
	DeleteTwoFA(ctx context.Context, userID string) error
}
