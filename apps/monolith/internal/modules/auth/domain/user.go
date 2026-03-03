package domain

import (
	"time"

	"gorm.io/gorm"
)

// User is the aggregate root for the auth bounded context.
type User struct {
	ID                       string     `gorm:"type:varchar(255);primaryKey" json:"id"`
	Email                    string     `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Password                 string     `gorm:"column:password_hash;size:255;not null" json:"-"`
	FirstName                string     `gorm:"size:100;not null" json:"first_name"`
	LastName                 string     `gorm:"size:100;not null" json:"last_name"`
	Phone                    string     `gorm:"size:20" json:"phone"`
	Avatar                   string     `gorm:"size:500" json:"avatar,omitempty"`
	IsActive                 bool       `gorm:"not null;default:true" json:"is_active"`
	IsVerified               bool       `gorm:"not null;default:false" json:"is_verified"`
	EmailVerificationToken   string     `gorm:"type:text" json:"-"`
	EmailVerificationExpires *time.Time `json:"-"`
	PasswordResetToken       string     `gorm:"type:text" json:"-"`
	PasswordResetExpires     *time.Time `json:"-"`
	LastLogin                *time.Time `json:"last_login,omitempty"`
	FailedLoginAttempts      int        `gorm:"not null;default:0" json:"-"`
	LockedUntil              *time.Time `json:"-"`
	CreatedAt                time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt                time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt                gorm.DeletedAt `gorm:"index" json:"-"`

	Preferences          *Preferences          `gorm:"foreignKey:UserID" json:"preferences,omitempty"`
	NotificationSettings *NotificationSettings `gorm:"foreignKey:UserID" json:"notification_settings,omitempty"`
	TwoFA                *TwoFA                `gorm:"foreignKey:UserID" json:"two_fa,omitempty"`
}

// TableName overrides the default GORM table name.
func (User) TableName() string { return "users" }

// ToResponse converts a User entity to a UserResponse DTO.
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:         u.ID,
		Email:      u.Email,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		Phone:      u.Phone,
		Avatar:     u.Avatar,
		IsActive:   u.IsActive,
		IsVerified: u.IsVerified,
		LastLogin:  u.LastLogin,
		CreatedAt:  u.CreatedAt,
	}
}
