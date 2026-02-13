package ports

import (
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/domain"
)

// JWTService handles JWT token generation and validation.
type JWTService interface {
	GenerateTokens(userID string, email string) (*domain.TokenPair, error)
	ValidateAccessToken(tokenString string) (*domain.Claims, error)
	ValidateRefreshToken(tokenString string) (*domain.Claims, error)
	GenerateEmailVerificationToken(userID string, email string) (string, error)
	GeneratePasswordResetToken(userID string, email string) (string, error)
	ValidateEmailVerificationToken(tokenString string) (*domain.Claims, error)
	ValidatePasswordResetToken(tokenString string) (*domain.Claims, error)
}

// PasswordService handles password hashing, verification, and strength validation.
type PasswordService interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) error
	ValidatePasswordStrength(password string) error
}

// TwoFAService handles TOTP two-factor authentication.
type TwoFAService interface {
	GenerateSecret(userEmail string) (secret string, qrCodeURL string, backupCodes []string, err error)
	GenerateQRCode(secret, userEmail string) ([]byte, error)
	ValidateCode(secret, code string) bool
	GenerateBackupCodes(count int) ([]string, error)
	ValidateBackupCode(codes []string, providedCode string) (remainingCodes []string, valid bool)
}
