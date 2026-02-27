package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/domain"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/ports"
	sharedports "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/shared/ports"
)

// Sentinel errors for auth flow control.
var (
	ErrTwoFARequired      = fmt.Errorf("2FA_REQUIRED")
	ErrInvalidCredentials  = fmt.Errorf("invalid email or password")
	ErrEmailAlreadyExists  = fmt.Errorf("user with email already exists")
	ErrAccountDeactivated  = fmt.Errorf("account is deactivated")
	ErrInvalid2FACode      = fmt.Errorf("invalid 2FA code")
)

// AuthService implements the core authentication business logic.
type AuthService struct {
	userRepo        ports.UserRepository
	prefsRepo       ports.PreferencesRepository
	notifsRepo      ports.NotificationSettingsRepository
	twoFARepo       ports.TwoFARepository
	jwtService      ports.JWTService
	passwordService ports.PasswordService
	twoFAService    ports.TwoFAService
	tenantCreator   ports.TenantCreator
	tenantFinder    ports.TenantMemberFinder
	eventBus        sharedports.EventBus
	logger          zerolog.Logger

	maxLoginAttempts int
	lockoutDuration  time.Duration
}

// NewAuthService creates a new AuthService with all required dependencies.
func NewAuthService(
	userRepo ports.UserRepository,
	prefsRepo ports.PreferencesRepository,
	notifsRepo ports.NotificationSettingsRepository,
	twoFARepo ports.TwoFARepository,
	jwtSvc ports.JWTService,
	pwSvc ports.PasswordService,
	twoFASvc ports.TwoFAService,
	tenantCreator ports.TenantCreator,
	tenantFinder ports.TenantMemberFinder,
	eventBus sharedports.EventBus,
	logger zerolog.Logger,
	maxLoginAttempts int,
	lockoutDuration time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		prefsRepo:        prefsRepo,
		notifsRepo:       notifsRepo,
		twoFARepo:        twoFARepo,
		jwtService:       jwtSvc,
		passwordService:  pwSvc,
		twoFAService:     twoFASvc,
		tenantCreator:    tenantCreator,
		tenantFinder:     tenantFinder,
		eventBus:         eventBus,
		logger:           logger,
		maxLoginAttempts: maxLoginAttempts,
		lockoutDuration:  lockoutDuration,
	}
}

// Register creates a new user account, sets up default preferences and
// notification settings, generates an email verification token, issues JWT
// tokens, and publishes a UserRegisteredEvent.
func (s *AuthService) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	// Check email uniqueness.
	existing, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil && existing != nil {
		s.logger.Warn().
			Str("component", "auth").
			Str("email", req.Email).
			Msg("registration attempt with existing email")
		return nil, ErrEmailAlreadyExists
	}

	// Hash the password.
	hashedPassword, err := s.passwordService.HashPassword(req.Password)
	if err != nil {
		s.logger.Error().
			Str("component", "auth").
			Err(err).
			Msg("failed to hash password during registration")
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Build the user entity.
	user := &domain.User{
		ID:         uuid.New().String(),
		Email:      req.Email,
		Password:   hashedPassword,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Phone:      req.Phone,
		IsActive:   true,
		IsVerified: false,
	}

	// Persist the user.
	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error().
			Str("component", "auth").
			Err(err).
			Msg("failed to create user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create default preferences.
	if _, err := s.prefsRepo.CreateDefaultPreferences(ctx, user.ID); err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", user.ID).
			Err(err).
			Msg("failed to create default preferences")
		return nil, fmt.Errorf("failed to create default preferences: %w", err)
	}

	// Create default notification settings.
	if _, err := s.notifsRepo.CreateDefaultNotifications(ctx, user.ID); err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", user.ID).
			Err(err).
			Msg("failed to create default notification settings")
		return nil, fmt.Errorf("failed to create default notification settings: %w", err)
	}

	// Generate and persist the email verification token.
	verificationToken, err := s.jwtService.GenerateEmailVerificationToken(user.ID, user.Email)
	if err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", user.ID).
			Err(err).
			Msg("failed to generate email verification token")
		return nil, fmt.Errorf("failed to generate email verification token: %w", err)
	}

	expires := time.Now().UTC().Add(24 * time.Hour)
	user.EmailVerificationToken = verificationToken
	user.EmailVerificationExpires = &expires

	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", user.ID).
			Err(err).
			Msg("failed to save email verification token")
		return nil, fmt.Errorf("failed to save email verification token: %w", err)
	}

	// Create personal tenant for the new user
	tenantID, tenantErr := s.tenantCreator.CreatePersonalTenant(ctx, user.ID, user.Email)
	if tenantErr != nil {
		s.logger.Error().Err(tenantErr).Str("user_id", user.ID).Msg("failed to create personal tenant")
		return nil, fmt.Errorf("failed to initialize tenant: %w", tenantErr)
	}

	// Generate JWT token pair.
	tokens, err := s.jwtService.GenerateTokens(user.ID, user.Email, tenantID, "owner")
	if err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", user.ID).
			Err(err).
			Msg("failed to generate JWT tokens")
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Publish the user registered event.
	event := domain.NewUserRegisteredEvent(user.ID, user.Email)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", user.ID).
			Err(err).
			Msg("failed to publish user registered event")
		// Non-fatal: continue despite event publishing failure.
	}

	s.logger.Info().
		Str("component", "auth").
		Str("user_id", user.ID).
		Str("email", user.Email).
		Msg("user registered successfully")

	return &domain.AuthResponse{
		User:   user.ToResponse(),
		Tokens: *tokens,
	}, nil
}

// Login authenticates a user by email and password, enforces account lockout
// policies, validates optional 2FA, issues JWT tokens, and publishes a
// UserLoggedInEvent.
func (s *AuthService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
	// Find user by email.
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil || user == nil {
		s.logger.Warn().
			Str("component", "auth").
			Str("email", req.Email).
			Msg("login attempt with unknown email")
		return nil, ErrInvalidCredentials
	}

	// Check if account is active.
	if !user.IsActive {
		s.logger.Warn().
			Str("component", "auth").
			Str("user_id", user.ID).
			Msg("login attempt on deactivated account")
		return nil, ErrAccountDeactivated
	}

	// Check if account is locked.
	if user.LockedUntil != nil && user.LockedUntil.After(time.Now().UTC()) {
		s.logger.Warn().
			Str("component", "auth").
			Str("user_id", user.ID).
			Time("locked_until", *user.LockedUntil).
			Msg("login attempt on locked account")
		return nil, fmt.Errorf("account is locked until %s", user.LockedUntil.Format(time.RFC3339))
	}

	// Verify password.
	if err := s.passwordService.VerifyPassword(user.Password, req.Password); err != nil {
		s.logger.Warn().
			Str("component", "auth").
			Str("user_id", user.ID).
			Int("failed_attempts", user.FailedLoginAttempts+1).
			Msg("invalid password attempt")

		user.FailedLoginAttempts++

		if user.FailedLoginAttempts >= s.maxLoginAttempts {
			lockUntil := time.Now().UTC().Add(s.lockoutDuration)
			user.LockedUntil = &lockUntil
			s.logger.Warn().
				Str("component", "auth").
				Str("user_id", user.ID).
				Time("locked_until", lockUntil).
				Msg("account locked due to too many failed attempts")
		}

		if updateErr := s.userRepo.Update(ctx, user); updateErr != nil {
			s.logger.Error().
				Str("component", "auth").
				Str("user_id", user.ID).
				Err(updateErr).
				Msg("failed to update failed login attempts")
		}

		return nil, ErrInvalidCredentials
	}

	// Check 2FA requirement.
	twoFA, err := s.twoFARepo.FindTwoFAByUserID(ctx, user.ID)
	if err == nil && twoFA != nil && twoFA.Enabled {
		if req.TwoFACode == "" {
			// 2FA is enabled but no code provided — client must prompt for code.
			return nil, ErrTwoFARequired
		}

		// Try TOTP validation first.
		if !s.twoFAService.ValidateCode(twoFA.Secret, req.TwoFACode) {
			// Fall back to backup code validation.
			remainingCodes, valid := s.twoFAService.ValidateBackupCode(twoFA.BackupCodes, req.TwoFACode)
			if !valid {
				s.logger.Warn().
					Str("component", "auth").
					Str("user_id", user.ID).
					Msg("invalid 2FA code during login")
				return nil, ErrInvalid2FACode
			}

			// Update the remaining backup codes.
			twoFA.BackupCodes = remainingCodes
			if upsertErr := s.twoFARepo.UpsertTwoFA(ctx, twoFA); upsertErr != nil {
				s.logger.Error().
					Str("component", "auth").
					Str("user_id", user.ID).
					Err(upsertErr).
					Msg("failed to update backup codes after use")
			}
		}
	}

	// Successful login — reset counters and record last login time.
	user.FailedLoginAttempts = 0
	user.LockedUntil = nil
	now := time.Now().UTC()
	user.LastLogin = &now

	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", user.ID).
			Err(err).
			Msg("failed to update user after successful login")
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Load tenant context for the authenticated user
	tenantID, role, tenantErr := s.tenantFinder.FindTenantByUserID(ctx, user.ID)
	if tenantErr != nil {
		s.logger.Error().Err(tenantErr).Str("user_id", user.ID).Msg("failed to load tenant context")
		return nil, fmt.Errorf("failed to load tenant context: %w", tenantErr)
	}

	// Generate JWT token pair.
	tokens, err := s.jwtService.GenerateTokens(user.ID, user.Email, tenantID, role)
	if err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", user.ID).
			Err(err).
			Msg("failed to generate JWT tokens during login")
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Publish the user logged in event.
	event := domain.NewUserLoggedInEvent(user.ID)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", user.ID).
			Err(err).
			Msg("failed to publish user logged in event")
		// Non-fatal: continue despite event publishing failure.
	}

	s.logger.Info().
		Str("component", "auth").
		Str("user_id", user.ID).
		Str("email", user.Email).
		Msg("user logged in successfully")

	return &domain.AuthResponse{
		User:   user.ToResponse(),
		Tokens: *tokens,
	}, nil
}

// Check2FA looks up whether a user (identified by email) has two-factor
// authentication enabled.
func (s *AuthService) Check2FA(ctx context.Context, email string) (*domain.Check2FAResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		s.logger.Warn().
			Str("component", "auth").
			Str("email", email).
			Msg("2FA check for unknown email")
		return nil, ErrInvalidCredentials
	}

	twoFA, err := s.twoFARepo.FindTwoFAByUserID(ctx, user.ID)
	if err != nil || twoFA == nil {
		return &domain.Check2FAResponse{TwoFAEnabled: false}, nil
	}

	return &domain.Check2FAResponse{TwoFAEnabled: twoFA.Enabled}, nil
}

// ---------------------------------------------------------------------------
// US2: Password & 2FA Management
// ---------------------------------------------------------------------------

// ChangePassword verifies the user's current password, validates and hashes the
// new one, persists the change, and publishes a PasswordChangedEvent.
func (s *AuthService) ChangePassword(ctx context.Context, userID string, req *domain.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", userID).
			Err(err).
			Msg("failed to find user for password change")
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Verify the current password.
	if err := s.passwordService.VerifyPassword(user.Password, req.CurrentPassword); err != nil {
		s.logger.Warn().
			Str("component", "auth").
			Str("user_id", userID).
			Msg("invalid current password during password change")
		return ErrInvalidCredentials
	}

	// Validate the strength of the new password.
	if err := s.passwordService.ValidatePasswordStrength(req.NewPassword); err != nil {
		return fmt.Errorf("password validation failed: %w", err)
	}

	// Hash the new password.
	hashedPassword, err := s.passwordService.HashPassword(req.NewPassword)
	if err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", userID).
			Err(err).
			Msg("failed to hash new password")
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", userID).
			Err(err).
			Msg("failed to update password")
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Publish password-changed event.
	event := domain.NewUserPasswordChangedEvent(userID)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", userID).
			Err(err).
			Msg("failed to publish password changed event")
	}

	s.logger.Info().
		Str("component", "auth").
		Str("user_id", userID).
		Msg("password changed successfully")

	return nil
}

// RequestPasswordReset generates a reset token and saves it on the user record.
// For security, the method returns nil even when the email is not found so that
// callers cannot enumerate valid accounts.
func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		// Return nil to avoid email enumeration attacks.
		s.logger.Info().
			Str("component", "auth").
			Str("email", email).
			Msg("password reset requested for unknown email — suppressed")
		return nil
	}

	// Generate a password-reset JWT.
	resetToken, err := s.jwtService.GeneratePasswordResetToken(user.ID, user.Email)
	if err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", user.ID).
			Err(err).
			Msg("failed to generate password reset token")
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	expires := time.Now().UTC().Add(1 * time.Hour)
	user.PasswordResetToken = resetToken
	user.PasswordResetExpires = &expires

	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", user.ID).
			Err(err).
			Msg("failed to save password reset token")
		return fmt.Errorf("failed to save reset token: %w", err)
	}

	s.logger.Info().
		Str("component", "auth").
		Str("user_id", user.ID).
		Msg("password reset token generated")

	return nil
}

// ResetPassword validates the reset JWT, hashes the new password, clears the
// reset token, resets lockout counters, and persists the user.
func (s *AuthService) ResetPassword(ctx context.Context, token string, newPassword string) error {
	// Validate the reset JWT.
	claims, err := s.jwtService.ValidatePasswordResetToken(token)
	if err != nil {
		s.logger.Warn().
			Str("component", "auth").
			Err(err).
			Msg("invalid password reset token")
		return fmt.Errorf("invalid or expired reset token: %w", err)
	}

	// Find the user by the stored reset token.
	user, err := s.userRepo.FindByResetToken(ctx, token)
	if err != nil || user == nil {
		s.logger.Warn().
			Str("component", "auth").
			Msg("no user found for reset token")
		return fmt.Errorf("invalid or expired reset token")
	}

	// Ensure the JWT belongs to the same user.
	if claims.UserID != user.ID {
		s.logger.Warn().
			Str("component", "auth").
			Str("claims_user_id", claims.UserID).
			Str("token_user_id", user.ID).
			Msg("reset token user ID mismatch")
		return fmt.Errorf("invalid or expired reset token")
	}

	// Validate strength.
	if err := s.passwordService.ValidatePasswordStrength(newPassword); err != nil {
		return fmt.Errorf("password validation failed: %w", err)
	}

	// Hash the new password.
	hashedPassword, err := s.passwordService.HashPassword(newPassword)
	if err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", user.ID).
			Err(err).
			Msg("failed to hash new password during reset")
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = hashedPassword
	user.PasswordResetToken = ""
	user.PasswordResetExpires = nil
	user.FailedLoginAttempts = 0
	user.LockedUntil = nil

	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", user.ID).
			Err(err).
			Msg("failed to update user after password reset")
		return fmt.Errorf("failed to update user: %w", err)
	}

	s.logger.Info().
		Str("component", "auth").
		Str("user_id", user.ID).
		Msg("password reset completed successfully")

	return nil
}

// Setup2FA generates a TOTP secret and QR code for the user and stores the
// TwoFA record in a non-enabled state until the user confirms with Enable2FA.
func (s *AuthService) Setup2FA(ctx context.Context, userID string) (*domain.TwoFASetupResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", userID).
			Err(err).
			Msg("failed to find user for 2FA setup")
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	secret, qrBase64, backupCodes, err := s.twoFAService.GenerateSecret(user.Email)
	if err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", userID).
			Err(err).
			Msg("failed to generate 2FA secret")
		return nil, fmt.Errorf("failed to generate 2FA secret: %w", err)
	}

	twoFA := &domain.TwoFA{
		UserID:      userID,
		Secret:      secret,
		Enabled:     false,
		BackupCodes: backupCodes,
	}

	if err := s.twoFARepo.UpsertTwoFA(ctx, twoFA); err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", userID).
			Err(err).
			Msg("failed to save 2FA record during setup")
		return nil, fmt.Errorf("failed to save 2FA record: %w", err)
	}

	s.logger.Info().
		Str("component", "auth").
		Str("user_id", userID).
		Msg("2FA setup initiated")

	return &domain.TwoFASetupResponse{
		Secret:      secret,
		QRCode:      qrBase64,
		BackupCodes: backupCodes,
	}, nil
}

// Enable2FA activates two-factor authentication after the user has verified the
// TOTP code from their authenticator app.
func (s *AuthService) Enable2FA(ctx context.Context, userID string, code string) error {
	twoFA, err := s.twoFARepo.FindTwoFAByUserID(ctx, userID)
	if err != nil || twoFA == nil {
		s.logger.Warn().
			Str("component", "auth").
			Str("user_id", userID).
			Msg("2FA record not found during enable")
		return fmt.Errorf("2FA not set up — call Setup2FA first")
	}

	if !s.twoFAService.ValidateCode(twoFA.Secret, code) {
		s.logger.Warn().
			Str("component", "auth").
			Str("user_id", userID).
			Msg("invalid 2FA code during enable")
		return ErrInvalid2FACode
	}

	twoFA.Enabled = true
	if err := s.twoFARepo.UpsertTwoFA(ctx, twoFA); err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", userID).
			Err(err).
			Msg("failed to enable 2FA")
		return fmt.Errorf("failed to enable 2FA: %w", err)
	}

	// Publish 2FA-enabled event.
	event := domain.NewUser2FAEnabledEvent(userID)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", userID).
			Err(err).
			Msg("failed to publish 2FA enabled event")
	}

	s.logger.Info().
		Str("component", "auth").
		Str("user_id", userID).
		Msg("2FA enabled successfully")

	return nil
}

// Disable2FA removes two-factor authentication after verifying the user's
// password.
func (s *AuthService) Disable2FA(ctx context.Context, userID string, password string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", userID).
			Err(err).
			Msg("failed to find user for 2FA disable")
		return fmt.Errorf("failed to find user: %w", err)
	}

	if err := s.passwordService.VerifyPassword(user.Password, password); err != nil {
		s.logger.Warn().
			Str("component", "auth").
			Str("user_id", userID).
			Msg("invalid password during 2FA disable")
		return ErrInvalidCredentials
	}

	if err := s.twoFARepo.DeleteTwoFA(ctx, userID); err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", userID).
			Err(err).
			Msg("failed to delete 2FA record")
		return fmt.Errorf("failed to disable 2FA: %w", err)
	}

	s.logger.Info().
		Str("component", "auth").
		Str("user_id", userID).
		Msg("2FA disabled successfully")

	return nil
}

// Verify2FA validates a TOTP code (or backup code) for an authenticated user,
// updating the remaining backup codes if one was consumed.
func (s *AuthService) Verify2FA(ctx context.Context, userID string, code string) error {
	twoFA, err := s.twoFARepo.FindTwoFAByUserID(ctx, userID)
	if err != nil || twoFA == nil {
		s.logger.Warn().
			Str("component", "auth").
			Str("user_id", userID).
			Msg("2FA record not found during verification")
		return fmt.Errorf("2FA is not set up")
	}

	if !twoFA.Enabled {
		return fmt.Errorf("2FA is not enabled")
	}

	// Try TOTP validation first.
	if s.twoFAService.ValidateCode(twoFA.Secret, code) {
		s.logger.Info().
			Str("component", "auth").
			Str("user_id", userID).
			Msg("2FA code verified via TOTP")
		return nil
	}

	// Fall back to backup code.
	remainingCodes, valid := s.twoFAService.ValidateBackupCode(twoFA.BackupCodes, code)
	if !valid {
		s.logger.Warn().
			Str("component", "auth").
			Str("user_id", userID).
			Msg("invalid 2FA code during verification")
		return ErrInvalid2FACode
	}

	// Update the remaining backup codes.
	twoFA.BackupCodes = remainingCodes
	if err := s.twoFARepo.UpsertTwoFA(ctx, twoFA); err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", userID).
			Err(err).
			Msg("failed to update backup codes after verification")
	}

	s.logger.Info().
		Str("component", "auth").
		Str("user_id", userID).
		Msg("2FA code verified via backup code")

	return nil
}

// ---------------------------------------------------------------------------
// US3: Session Management
// ---------------------------------------------------------------------------

// RefreshToken validates the provided refresh JWT and issues a new token pair.
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	claims, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		s.logger.Warn().
			Str("component", "auth").
			Err(err).
			Msg("invalid refresh token")
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil || user == nil {
		s.logger.Warn().
			Str("component", "auth").
			Str("user_id", claims.UserID).
			Msg("user not found during token refresh")
		return nil, fmt.Errorf("user not found")
	}

	if !user.IsActive {
		s.logger.Warn().
			Str("component", "auth").
			Str("user_id", user.ID).
			Msg("token refresh attempt on deactivated account")
		return nil, ErrAccountDeactivated
	}

	tokens, err := s.jwtService.GenerateTokens(claims.UserID, claims.Email, claims.TenantID, claims.Role)
	if err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", user.ID).
			Err(err).
			Msg("failed to generate tokens during refresh")
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	s.logger.Info().
		Str("component", "auth").
		Str("user_id", user.ID).
		Msg("tokens refreshed successfully")

	return tokens, nil
}

// Logout is a no-op for stateless JWT authentication. It exists to satisfy the
// service interface and allow future implementations (e.g., token blacklisting).
func (s *AuthService) Logout(ctx context.Context, userID string) error {
	s.logger.Info().
		Str("component", "auth").
		Str("user_id", userID).
		Msg("user logged out (stateless — no-op)")
	return nil
}

// ---------------------------------------------------------------------------
// US4: Profile & Preferences
// ---------------------------------------------------------------------------

// GetProfile returns the public profile representation for the given user.
func (s *AuthService) GetProfile(ctx context.Context, userID string) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", userID).
			Err(err).
			Msg("failed to find user for profile retrieval")
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	resp := user.ToResponse()
	return &resp, nil
}

// UpdateProfile updates the user's name and phone fields and returns the
// updated public profile. Avatar is handled separately via UploadAvatar.
func (s *AuthService) UpdateProfile(ctx context.Context, userID string, req *domain.UpdateProfileRequest) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", userID).
			Err(err).
			Msg("failed to find user for profile update")
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Phone = req.Phone

	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", userID).
			Err(err).
			Msg("failed to update user profile")
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	s.logger.Info().
		Str("component", "auth").
		Str("user_id", userID).
		Msg("profile updated successfully")

	resp := user.ToResponse()
	return &resp, nil
}

// UploadAvatar sets the avatar path for the given user.
func (s *AuthService) UploadAvatar(ctx context.Context, userID string, avatarPath string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", userID).
			Err(err).
			Msg("failed to find user for avatar upload")
		return fmt.Errorf("failed to find user: %w", err)
	}

	user.Avatar = avatarPath

	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", userID).
			Err(err).
			Msg("failed to update avatar")
		return fmt.Errorf("failed to update avatar: %w", err)
	}

	s.logger.Info().
		Str("component", "auth").
		Str("user_id", userID).
		Str("avatar_path", avatarPath).
		Msg("avatar updated successfully")

	return nil
}

// GetPreferences returns the user's display and locale preferences.
func (s *AuthService) GetPreferences(ctx context.Context, userID string) (*domain.Preferences, error) {
	prefs, err := s.prefsRepo.FindPreferencesByUserID(ctx, userID)
	if err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", userID).
			Err(err).
			Msg("failed to retrieve preferences")
		return nil, fmt.Errorf("failed to get preferences: %w", err)
	}
	return prefs, nil
}

// UpdatePreferences updates the user's display and locale preferences.
func (s *AuthService) UpdatePreferences(ctx context.Context, userID string, prefs *domain.Preferences) error {
	prefs.UserID = userID

	if err := s.prefsRepo.UpdatePreferences(ctx, prefs); err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", userID).
			Err(err).
			Msg("failed to update preferences")
		return fmt.Errorf("failed to update preferences: %w", err)
	}

	s.logger.Info().
		Str("component", "auth").
		Str("user_id", userID).
		Msg("preferences updated successfully")

	return nil
}

// GetNotifications returns the user's notification settings.
func (s *AuthService) GetNotifications(ctx context.Context, userID string) (*domain.NotificationSettings, error) {
	settings, err := s.notifsRepo.FindNotificationsByUserID(ctx, userID)
	if err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", userID).
			Err(err).
			Msg("failed to retrieve notification settings")
		return nil, fmt.Errorf("failed to get notification settings: %w", err)
	}
	return settings, nil
}

// UpdateNotifications updates the user's notification settings.
func (s *AuthService) UpdateNotifications(ctx context.Context, userID string, settings *domain.NotificationSettings) error {
	settings.UserID = userID

	if err := s.notifsRepo.UpdateNotifications(ctx, settings); err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", userID).
			Err(err).
			Msg("failed to update notification settings")
		return fmt.Errorf("failed to update notification settings: %w", err)
	}

	s.logger.Info().
		Str("component", "auth").
		Str("user_id", userID).
		Msg("notification settings updated successfully")

	return nil
}

// ---------------------------------------------------------------------------
// US5: Email Verification
// ---------------------------------------------------------------------------

// VerifyEmail validates the email verification JWT, locates the user by the
// stored verification token, marks the account as verified, and clears the
// verification fields.
func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	claims, err := s.jwtService.ValidateEmailVerificationToken(token)
	if err != nil {
		s.logger.Warn().
			Str("component", "auth").
			Err(err).
			Msg("invalid email verification token")
		return fmt.Errorf("invalid or expired verification token: %w", err)
	}

	user, err := s.userRepo.FindByVerificationToken(ctx, token)
	if err != nil || user == nil {
		s.logger.Warn().
			Str("component", "auth").
			Msg("no user found for verification token")
		return fmt.Errorf("invalid or expired verification token")
	}

	if claims.UserID != user.ID {
		s.logger.Warn().
			Str("component", "auth").
			Str("claims_user_id", claims.UserID).
			Str("token_user_id", user.ID).
			Msg("verification token user ID mismatch")
		return fmt.Errorf("invalid or expired verification token")
	}

	user.IsVerified = true
	user.EmailVerificationToken = ""
	user.EmailVerificationExpires = nil

	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", user.ID).
			Err(err).
			Msg("failed to update user after email verification")
		return fmt.Errorf("failed to verify email: %w", err)
	}

	s.logger.Info().
		Str("component", "auth").
		Str("user_id", user.ID).
		Msg("email verified successfully")

	return nil
}

// ---------------------------------------------------------------------------
// US6: Account Data Management
// ---------------------------------------------------------------------------

// ExportData gathers the user's profile, preferences, and notification settings
// into a single map suitable for JSON serialization and download.
func (s *AuthService) ExportData(ctx context.Context, userID string) (map[string]interface{}, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", userID).
			Err(err).
			Msg("failed to find user for data export")
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	prefs, err := s.prefsRepo.FindPreferencesByUserID(ctx, userID)
	if err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", userID).
			Err(err).
			Msg("failed to retrieve preferences for data export")
		return nil, fmt.Errorf("failed to get preferences: %w", err)
	}

	notifications, err := s.notifsRepo.FindNotificationsByUserID(ctx, userID)
	if err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", userID).
			Err(err).
			Msg("failed to retrieve notifications for data export")
		return nil, fmt.Errorf("failed to get notification settings: %w", err)
	}

	data := map[string]interface{}{
		"profile":       user.ToResponse(),
		"preferences":   prefs,
		"notifications": notifications,
	}

	s.logger.Info().
		Str("component", "auth").
		Str("user_id", userID).
		Msg("user data exported successfully")

	return data, nil
}

// DeleteAccount verifies the user's password, permanently deletes the account,
// and publishes a UserDeletedEvent.
func (s *AuthService) DeleteAccount(ctx context.Context, userID string, password string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", userID).
			Err(err).
			Msg("failed to find user for account deletion")
		return fmt.Errorf("failed to find user: %w", err)
	}

	if err := s.passwordService.VerifyPassword(user.Password, password); err != nil {
		s.logger.Warn().
			Str("component", "auth").
			Str("user_id", userID).
			Msg("invalid password during account deletion")
		return ErrInvalidCredentials
	}

	if err := s.userRepo.Delete(ctx, userID); err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", userID).
			Err(err).
			Msg("failed to delete user account")
		return fmt.Errorf("failed to delete account: %w", err)
	}

	// Publish account-deleted event.
	event := domain.NewUserDeletedEvent(userID, user.Email)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		s.logger.Error().
			Str("component", "auth").
			Str("user_id", userID).
			Err(err).
			Msg("failed to publish user deleted event")
	}

	s.logger.Info().
		Str("component", "auth").
		Str("user_id", userID).
		Str("email", user.Email).
		Msg("account deleted successfully")

	return nil
}
