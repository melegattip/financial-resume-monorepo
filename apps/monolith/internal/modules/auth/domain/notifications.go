package domain

import (
	"time"

	"gorm.io/gorm"
)

// NotificationSettings holds per-user notification preferences (1:1 with User).
type NotificationSettings struct {
	UserID             string         `gorm:"type:varchar(255);primaryKey" json:"user_id"`
	EmailNotifications bool           `gorm:"not null;default:false" json:"email_notifications"`
	BudgetAlerts       bool           `gorm:"not null;default:false" json:"budget_alerts"`
	LoginNotifications bool           `gorm:"not null;default:false" json:"login_notifications"`
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"-"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"-"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName overrides the default GORM table name.
func (NotificationSettings) TableName() string { return "user_notification_settings" }

// DefaultNotificationSettings returns settings with all notifications disabled (opt-in).
func DefaultNotificationSettings(userID string) NotificationSettings {
	return NotificationSettings{
		UserID:             userID,
		EmailNotifications: false,
		BudgetAlerts:       false,
		LoginNotifications: false,
	}
}
