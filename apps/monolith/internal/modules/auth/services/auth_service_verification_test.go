package services_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/domain"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/ports"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/services"
	sharedports "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/shared/ports"
)

// ---------------------------------------------------------------------------
// Minimal mock helpers
// ---------------------------------------------------------------------------

// stubUserRepo implements ports.UserRepository with function fields.
type stubUserRepo struct {
	findByEmail            func(ctx context.Context, email string) (*domain.User, error)
	findByID               func(ctx context.Context, id string) (*domain.User, error)
	findByVerificationToken func(ctx context.Context, token string) (*domain.User, error)
	findByResetToken       func(ctx context.Context, token string) (*domain.User, error)
	create                 func(ctx context.Context, user *domain.User) error
	update                 func(ctx context.Context, user *domain.User) error
	delete                 func(ctx context.Context, id string) error
}

func (s *stubUserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.findByEmail(ctx, email)
}
func (s *stubUserRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	if s.findByID != nil {
		return s.findByID(ctx, id)
	}
	return nil, fmt.Errorf("not found")
}
func (s *stubUserRepo) FindByVerificationToken(ctx context.Context, token string) (*domain.User, error) {
	if s.findByVerificationToken != nil {
		return s.findByVerificationToken(ctx, token)
	}
	return nil, fmt.Errorf("not found")
}
func (s *stubUserRepo) FindByResetToken(ctx context.Context, token string) (*domain.User, error) {
	if s.findByResetToken != nil {
		return s.findByResetToken(ctx, token)
	}
	return nil, fmt.Errorf("not found")
}
func (s *stubUserRepo) Create(ctx context.Context, user *domain.User) error {
	if s.create != nil {
		return s.create(ctx, user)
	}
	return nil
}
func (s *stubUserRepo) Update(ctx context.Context, user *domain.User) error {
	return s.update(ctx, user)
}
func (s *stubUserRepo) Delete(ctx context.Context, id string) error { return nil }

// noopPrefsRepo, noopNotifsRepo, noopTwoFARepo are empty stubs for unused repos.

type noopPrefsRepo struct{}

func (n *noopPrefsRepo) FindPreferencesByUserID(ctx context.Context, userID string) (*domain.Preferences, error) {
	return &domain.Preferences{}, nil
}
func (n *noopPrefsRepo) CreateDefaultPreferences(ctx context.Context, userID string) (*domain.Preferences, error) {
	return &domain.Preferences{}, nil
}
func (n *noopPrefsRepo) UpdatePreferences(ctx context.Context, prefs *domain.Preferences) error {
	return nil
}

type noopNotifsRepo struct{}

func (n *noopNotifsRepo) FindNotificationsByUserID(ctx context.Context, userID string) (*domain.NotificationSettings, error) {
	return &domain.NotificationSettings{}, nil
}
func (n *noopNotifsRepo) CreateDefaultNotifications(ctx context.Context, userID string) (*domain.NotificationSettings, error) {
	return &domain.NotificationSettings{}, nil
}
func (n *noopNotifsRepo) UpdateNotifications(ctx context.Context, settings *domain.NotificationSettings) error {
	return nil
}

type noopTwoFARepo struct{}

func (n *noopTwoFARepo) FindTwoFAByUserID(ctx context.Context, userID string) (*domain.TwoFA, error) {
	return nil, fmt.Errorf("not found")
}
func (n *noopTwoFARepo) UpsertTwoFA(ctx context.Context, twoFA *domain.TwoFA) error { return nil }
func (n *noopTwoFARepo) DeleteTwoFA(ctx context.Context, userID string) error        { return nil }

// stubJWTService captures generated tokens and returns predictable values.
type stubJWTService struct {
	generatedVerificationToken string
}

func (s *stubJWTService) GenerateTokens(userID, email, tenantID, role string) (*domain.TokenPair, error) {
	exp := time.Now().Add(time.Hour)
	return &domain.TokenPair{AccessToken: "access", RefreshToken: "refresh", ExpiresAt: exp, TokenType: "Bearer"}, nil
}
func (s *stubJWTService) ValidateAccessToken(token string) (*domain.Claims, error) {
	return nil, nil
}
func (s *stubJWTService) ValidateRefreshToken(token string) (*domain.Claims, error) {
	return nil, nil
}
func (s *stubJWTService) GenerateEmailVerificationToken(userID, email string) (string, error) {
	s.generatedVerificationToken = "stub-verification-token-for-" + userID
	return s.generatedVerificationToken, nil
}
func (s *stubJWTService) GeneratePasswordResetToken(userID, email string) (string, error) {
	return "stub-reset-token", nil
}
func (s *stubJWTService) ValidateEmailVerificationToken(token string) (*domain.Claims, error) {
	if token == "stub-verification-token-for-user-123" {
		return &domain.Claims{UserID: "user-123"}, nil
	}
	return nil, fmt.Errorf("invalid token")
}
func (s *stubJWTService) ValidatePasswordResetToken(token string) (*domain.Claims, error) {
	return nil, nil
}

// stubEmailSvc captures sent emails for assertion.
type stubEmailSvc struct {
	sentVerification []string // email addresses
	failWith         error
}

func (s *stubEmailSvc) SendEmailVerification(toEmail, firstName, link string) error {
	if s.failWith != nil {
		return s.failWith
	}
	s.sentVerification = append(s.sentVerification, toEmail)
	return nil
}
func (s *stubEmailSvc) SendPasswordReset(toEmail, resetLink string) error          { return nil }
func (s *stubEmailSvc) SendBudgetAlert(toEmail, firstName, categoryID, period, newStatus string, spentAmount, budgetLimit float64) error {
	return nil
}
func (s *stubEmailSvc) SendLoginNotification(toEmail, firstName, loginTime string) error { return nil }

// noopTenantCreator / Finder / Cleaner

type noopTenantCreator struct{}

func (n *noopTenantCreator) CreatePersonalTenant(ctx context.Context, userID, email string) (string, error) {
	return "tnt_test", nil
}

type noopTenantFinder struct{}

func (n *noopTenantFinder) FindTenantByUserID(ctx context.Context, userID string) (string, string, error) {
	return "tnt_test", "owner", nil
}
func (n *noopTenantFinder) FindMemberInTenant(ctx context.Context, userID, tenantID string) (string, error) {
	return "owner", nil
}

type noopTenantCleaner struct{}

func (n *noopTenantCleaner) CleanupUserTenants(ctx context.Context, userID string) error { return nil }

// noopEventBus

type noopEventBus struct{}

func (n *noopEventBus) Publish(ctx context.Context, event sharedports.Event) error { return nil }
func (n *noopEventBus) Subscribe(eventType string, handler sharedports.EventHandler) {}

// buildAuthService wires a minimal AuthService for testing.
func buildAuthService(userRepo ports.UserRepository, jwtSvc ports.JWTService, emailSvc ports.EmailService) *services.AuthService {
	return services.NewAuthService(
		userRepo,
		&noopPrefsRepo{},
		&noopNotifsRepo{},
		&noopTwoFARepo{},
		jwtSvc,
		services.NewPasswordService(8),
		services.NewTwoFAService("niloft"),
		&noopTenantCreator{},
		&noopTenantFinder{},
		&noopTenantCleaner{},
		emailSvc,
		"http://localhost:3000",
		&noopEventBus{},
		zerolog.Nop(),
		5,
		15*time.Minute,
	)
}

// ---------------------------------------------------------------------------
// ResendVerificationEmail tests
// ---------------------------------------------------------------------------

func TestResendVerificationEmail_SendsEmailForUnverifiedUser(t *testing.T) {
	unverifiedUser := &domain.User{
		ID:         "user-123",
		Email:      "user@example.com",
		FirstName:  "Test",
		IsVerified: false,
		IsActive:   true,
	}

	updateCalled := false
	repo := &stubUserRepo{
		findByEmail: func(_ context.Context, email string) (*domain.User, error) {
			return unverifiedUser, nil
		},
		update: func(_ context.Context, user *domain.User) error {
			updateCalled = true
			// Token should be set before the Update call
			assert.NotEmpty(t, user.EmailVerificationToken)
			assert.NotNil(t, user.EmailVerificationExpires)
			return nil
		},
	}

	jwtSvc := &stubJWTService{}
	emailSvc := &stubEmailSvc{}
	svc := buildAuthService(repo, jwtSvc, emailSvc)

	err := svc.ResendVerificationEmail(context.Background(), "user@example.com")
	require.NoError(t, err)
	assert.True(t, updateCalled, "Update must be called to persist the new token")

	// Give the goroutine time to send the email
	time.Sleep(50 * time.Millisecond)
	assert.Contains(t, emailSvc.sentVerification, "user@example.com")
}

func TestResendVerificationEmail_SilentForUnknownEmail(t *testing.T) {
	// Must return nil (no error) to prevent email enumeration.
	repo := &stubUserRepo{
		findByEmail: func(_ context.Context, email string) (*domain.User, error) {
			return nil, fmt.Errorf("user not found")
		},
		update: func(_ context.Context, _ *domain.User) error {
			t.Fatal("Update must not be called for unknown emails")
			return nil
		},
	}

	svc := buildAuthService(repo, &stubJWTService{}, &stubEmailSvc{})
	err := svc.ResendVerificationEmail(context.Background(), "nobody@example.com")
	assert.NoError(t, err)
}

func TestResendVerificationEmail_SilentForAlreadyVerifiedUser(t *testing.T) {
	verifiedUser := &domain.User{
		ID:         "user-999",
		Email:      "verified@example.com",
		IsVerified: true,
	}

	repo := &stubUserRepo{
		findByEmail: func(_ context.Context, _ string) (*domain.User, error) {
			return verifiedUser, nil
		},
		update: func(_ context.Context, _ *domain.User) error {
			t.Fatal("Update must not be called for already-verified users")
			return nil
		},
	}

	svc := buildAuthService(repo, &stubJWTService{}, &stubEmailSvc{})
	err := svc.ResendVerificationEmail(context.Background(), "verified@example.com")
	assert.NoError(t, err)
}

func TestResendVerificationEmail_TokenExpiresIn24Hours(t *testing.T) {
	unverifiedUser := &domain.User{ID: "user-123", Email: "user@example.com", IsVerified: false}

	var savedExpiry *time.Time
	repo := &stubUserRepo{
		findByEmail: func(_ context.Context, _ string) (*domain.User, error) { return unverifiedUser, nil },
		update: func(_ context.Context, user *domain.User) error {
			savedExpiry = user.EmailVerificationExpires
			return nil
		},
	}

	svc := buildAuthService(repo, &stubJWTService{}, &stubEmailSvc{})
	require.NoError(t, svc.ResendVerificationEmail(context.Background(), "user@example.com"))

	require.NotNil(t, savedExpiry)
	hoursUntilExpiry := time.Until(*savedExpiry).Hours()
	assert.InDelta(t, 24.0, hoursUntilExpiry, 0.1, "token must expire in ~24 hours")
}

func TestResendVerificationEmail_EmailLinkContainsToken(t *testing.T) {
	unverifiedUser := &domain.User{ID: "user-123", Email: "user@example.com", IsVerified: false}

	repo := &stubUserRepo{
		findByEmail: func(_ context.Context, _ string) (*domain.User, error) { return unverifiedUser, nil },
		update:      func(_ context.Context, _ *domain.User) error { return nil },
	}

	jwtSvc := &stubJWTService{}
	// Use a custom emailSvc that captures the link
	var capturedLink string
	emailSvc := &captureEmailSvc{onSend: func(_, _, link string) { capturedLink = link }}

	svc := buildAuthService(repo, jwtSvc, emailSvc)
	require.NoError(t, svc.ResendVerificationEmail(context.Background(), "user@example.com"))

	time.Sleep(50 * time.Millisecond)
	expectedToken := jwtSvc.generatedVerificationToken
	assert.Contains(t, capturedLink, "/verify-email?token="+expectedToken,
		"verification link must contain the generated token")
}

// ---------------------------------------------------------------------------
// VerifyEmail service tests
// ---------------------------------------------------------------------------

func TestVerifyEmail_MarksUserAsVerified(t *testing.T) {
	token := "stub-verification-token-for-user-123"
	exp := time.Now().Add(time.Hour)
	unverifiedUser := &domain.User{
		ID:                      "user-123",
		Email:                   "user@example.com",
		IsVerified:              false,
		EmailVerificationToken:  token,
		EmailVerificationExpires: &exp,
	}

	var savedUser *domain.User
	repo := &stubUserRepo{
		findByVerificationToken: func(_ context.Context, _ string) (*domain.User, error) {
			return unverifiedUser, nil
		},
		update: func(_ context.Context, user *domain.User) error {
			savedUser = user
			return nil
		},
	}

	svc := buildAuthService(repo, &stubJWTService{}, &stubEmailSvc{})
	err := svc.VerifyEmail(context.Background(), token)
	require.NoError(t, err)

	require.NotNil(t, savedUser)
	assert.True(t, savedUser.IsVerified)
	assert.Empty(t, savedUser.EmailVerificationToken, "token must be cleared after verification")
	assert.Nil(t, savedUser.EmailVerificationExpires, "expiry must be cleared after verification")
}

func TestVerifyEmail_RejectsInvalidToken(t *testing.T) {
	repo := &stubUserRepo{
		findByVerificationToken: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, fmt.Errorf("invalid or expired email verification token")
		},
		update: func(_ context.Context, _ *domain.User) error {
			t.Fatal("Update must not be called for invalid token")
			return nil
		},
	}

	svc := buildAuthService(repo, &stubJWTService{}, &stubEmailSvc{})
	err := svc.VerifyEmail(context.Background(), "completely-wrong-token")
	assert.Error(t, err)
}

// ---------------------------------------------------------------------------
// captureEmailSvc — captures link for assertion
// ---------------------------------------------------------------------------

type captureEmailSvc struct {
	onSend func(toEmail, firstName, link string)
}

func (c *captureEmailSvc) SendEmailVerification(toEmail, firstName, link string) error {
	c.onSend(toEmail, firstName, link)
	return nil
}
func (c *captureEmailSvc) SendPasswordReset(toEmail, resetLink string) error { return nil }
func (c *captureEmailSvc) SendBudgetAlert(toEmail, firstName, categoryID, period, newStatus string, spentAmount, budgetLimit float64) error {
	return nil
}
func (c *captureEmailSvc) SendLoginNotification(toEmail, firstName, loginTime string) error {
	return nil
}
