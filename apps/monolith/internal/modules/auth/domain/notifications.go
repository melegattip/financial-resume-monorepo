package domain

import (
	"time"

	"gorm.io/gorm"
)

// NotificationSettings holds per-user notification preferences (1:1 with User).
type NotificationSettings struct {
	UserID                   string    `gorm:"type:varchar(255);primaryKey" json:"user_id"`
	EmailNotifications       bool      `gorm:"not null;default:true" json:"email_notifications"`
	PushNotifications        bool      `gorm:"not null;default:true" json:"push_notifications"`
	WeeklyReports            bool      `gorm:"not null;default:true" json:"weekly_reports"`
	ExpenseAlerts            bool      `gorm:"not null;default:true" json:"expense_alerts"`
	BudgetAlerts             bool      `gorm:"not null;default:true" json:"budget_alerts"`
	AchievementNotifications bool      `gorm:"not null;default:true" json:"achievement_notifications"`
	CreatedAt                time.Time      `gorm:"autoCreateTime" json:"-"`
	UpdatedAt                time.Time      `gorm:"autoUpdateTime" json:"-"`
	DeletedAt                gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName overrides the default GORM table name.
func (NotificationSettings) TableName() string { return "user_notification_settings" }

// DefaultNotificationSettings returns settings with all notifications enabled.
func DefaultNotificationSettings(userID string) NotificationSettings {
	return NotificationSettings{
		UserID:                   userID,
		EmailNotifications:       true,
		PushNotifications:        true,
		WeeklyReports:            true,
		ExpenseAlerts:            true,
		BudgetAlerts:             true,
		AchievementNotifications: true,
	}
}
