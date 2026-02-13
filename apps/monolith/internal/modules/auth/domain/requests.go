package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// RegisterRequest is the DTO for user registration.
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone"`
}

// LoginRequest is the DTO for user login.
type LoginRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	TwoFACode string `json:"twofa_code,omitempty"`
}

// ChangePasswordRequest is the DTO for changing the current password.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// UpdateProfileRequest is the DTO for profile updates.
type UpdateProfileRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone"`
	Avatar    string `json:"avatar,omitempty"`
}

// PasswordResetRequest is the DTO for requesting a password reset.
type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest is the DTO for completing a password reset.
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// Enable2FARequest is the DTO for enabling 2FA after setup.
type Enable2FARequest struct {
	Code string `json:"code" binding:"required"`
}

// Verify2FARequest is the DTO for verifying a 2FA code.
type Verify2FARequest struct {
	Code string `json:"code" binding:"required"`
}

// Disable2FARequest is the DTO for disabling 2FA.
type Disable2FARequest struct {
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest is the DTO for refreshing an access token.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// DeleteAccountRequest is the DTO for account deletion confirmation.
type DeleteAccountRequest struct {
	Password string `json:"password" binding:"required"`
}

// UserResponse is the public representation of a user.
type UserResponse struct {
	ID         string     `json:"id"`
	Email      string     `json:"email"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Phone      string     `json:"phone"`
	Avatar     string     `json:"avatar,omitempty"`
	IsActive   bool       `json:"is_active"`
	IsVerified bool       `json:"is_verified"`
	LastLogin  *time.Time `json:"last_login,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// TokenPair holds the JWT access and refresh tokens returned to clients.
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

// Claims represents the custom JWT claims payload.
type Claims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// AuthResponse wraps tokens and user data for auth endpoint responses.
type AuthResponse struct {
	User   UserResponse `json:"user"`
	Tokens TokenPair    `json:"tokens"`
}

// TwoFASetupResponse is returned when 2FA is set up (before enabling).
type TwoFASetupResponse struct {
	Secret      string   `json:"secret"`
	QRCode      string   `json:"qr_code"`
	BackupCodes []string `json:"backup_codes"`
}

// Check2FAResponse indicates whether a user has 2FA enabled.
type Check2FAResponse struct {
	TwoFAEnabled bool `json:"twofa_enabled"`
}
