package domain

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// TwoFA holds TOTP two-factor authentication configuration (1:0..1 with User).
type TwoFA struct {
	UserID       string         `gorm:"type:varchar(255);primaryKey" json:"user_id"`
	Secret       string         `gorm:"size:255;not null" json:"-"`
	Enabled      bool           `gorm:"not null;default:false" json:"enabled"`
	BackupCodes  pq.StringArray `gorm:"type:text[]" json:"backup_codes,omitempty"`
	LastUsedCode string         `gorm:"size:50" json:"-"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"-"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"-"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName overrides the default GORM table name.
func (TwoFA) TableName() string { return "user_two_fa" }
