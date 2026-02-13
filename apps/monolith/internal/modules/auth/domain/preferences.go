package domain

import (
	"time"

	"gorm.io/gorm"
)

// Preferences holds user display and locale settings (1:1 with User).
type Preferences struct {
	UserID     string    `gorm:"type:varchar(255);primaryKey" json:"user_id"`
	Currency   string    `gorm:"size:10;not null;default:'USD'" json:"currency"`
	Language   string    `gorm:"size:10;not null;default:'en'" json:"language"`
	Theme      string    `gorm:"size:10;not null;default:'light'" json:"theme"`
	DateFormat string    `gorm:"size:20;not null;default:'YYYY-MM-DD'" json:"date_format"`
	Timezone   string    `gorm:"size:50;not null;default:'UTC'" json:"timezone"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"-"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName overrides the default GORM table name.
func (Preferences) TableName() string { return "user_preferences" }

// DefaultPreferences returns a new Preferences with default values for the given user.
func DefaultPreferences(userID string) Preferences {
	return Preferences{
		UserID:     userID,
		Currency:   "USD",
		Language:   "en",
		Theme:      "light",
		DateFormat: "YYYY-MM-DD",
		Timezone:   "UTC",
	}
}
