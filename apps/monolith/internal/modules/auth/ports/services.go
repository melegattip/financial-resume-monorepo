package ports

import (
	"context"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/domain"
)

// JWTService handles JWT token generation and validation.
type JWTService interface {
	GenerateTokens(userID string, email string, tenantID string, role string) (*domain.TokenPair, error)
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

// EmailService handles sending transactional emails.
type EmailService interface {
	SendPasswordReset(toEmail, resetLink string) error
	SendEmailVerification(toEmail, firstName, verificationLink string) error
	SendBudgetAlert(toEmail, firstName, categoryID, period, newStatus string, spentAmount, budgetLimit float64) error
	SendLoginNotification(toEmail, firstName, loginTime string) error
}

// TenantCreator allows the auth module to create a personal tenant for new users
// without importing the tenants module (avoids circular imports).
type TenantCreator interface {
	CreatePersonalTenant(ctx context.Context, userID string, email string) (tenantID string, err error)
}

// TenantMemberFinder allows the auth module to load tenant context on login and switch tenants.
type TenantMemberFinder interface {
	FindTenantByUserID(ctx context.Context, userID string) (tenantID string, role string, err error)
	// FindMemberInTenant returns the user's role in a specific tenant.
	// Returns an empty string and nil error when the user is not a member.
	FindMemberInTenant(ctx context.Context, userID, tenantID string) (role string, err error)
}

// TenantAccountCleaner handles tenant cleanup when a user account is deleted.
// For each tenant the user owns: if other members exist the next eligible member
// (highest role then oldest join date) becomes the new owner; otherwise the
// tenant is soft-deleted. The user is then removed from all memberships.
type TenantAccountCleaner interface {
	CleanupUserTenants(ctx context.Context, userID string) error
}
