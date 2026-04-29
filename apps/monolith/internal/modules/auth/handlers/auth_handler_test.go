package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/domain"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/handlers"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// ---------------------------------------------------------------------------
// Mock implementation of authSvc
// ---------------------------------------------------------------------------

type mockAuthSvc struct {
	registerFn              func(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error)
	loginFn                 func(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error)
	check2FAFn              func(ctx context.Context, email string) (*domain.Check2FAResponse, error)
	switchTenantFn          func(ctx context.Context, userID, tenantID string) (*domain.TokenPair, error)
	verifyEmailFn           func(ctx context.Context, token string) error
	resendVerificationFn    func(ctx context.Context, email string) error
}

func (m *mockAuthSvc) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	return m.registerFn(ctx, req)
}
func (m *mockAuthSvc) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
	return m.loginFn(ctx, req)
}
func (m *mockAuthSvc) Check2FA(ctx context.Context, email string) (*domain.Check2FAResponse, error) {
	return m.check2FAFn(ctx, email)
}
func (m *mockAuthSvc) SwitchTenant(ctx context.Context, userID, tenantID string) (*domain.TokenPair, error) {
	return m.switchTenantFn(ctx, userID, tenantID)
}
func (m *mockAuthSvc) VerifyEmail(ctx context.Context, token string) error {
	return m.verifyEmailFn(ctx, token)
}
func (m *mockAuthSvc) ResendVerificationEmail(ctx context.Context, email string) error {
	return m.resendVerificationFn(ctx, email)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func newTestRouter(h *handlers.AuthHandler) *gin.Engine {
	r := gin.New()
	r.POST("/auth/register", h.Register)
	r.POST("/auth/login", h.Login)
	r.GET("/auth/verify-email/:token", h.VerifyEmail)
	r.POST("/auth/resend-verification", h.ResendVerification)
	return r
}

func doRequest(r *gin.Engine, method, path string, body any) *httptest.ResponseRecorder {
	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = json.Marshal(body)
	}
	req := httptest.NewRequest(method, path, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func parseBody(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var out map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &out))
	return out
}

func silentLogger() zerolog.Logger { return zerolog.Nop() }

// ---------------------------------------------------------------------------
// VerifyEmail handler tests
// ---------------------------------------------------------------------------

func TestVerifyEmail_Success(t *testing.T) {
	svc := &mockAuthSvc{verifyEmailFn: func(_ context.Context, token string) error { return nil }}
	r := newTestRouter(handlers.NewAuthHandler(svc, silentLogger()))

	w := doRequest(r, http.MethodGet, "/auth/verify-email/valid-token-123", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	body := parseBody(t, w)
	assert.Equal(t, "Email verified successfully", body["message"])
}

func TestVerifyEmail_MissingToken_NotFound(t *testing.T) {
	// Route requires /:token — hitting the base path without a param segment
	// returns a non-200 status from the router (404 in test mode, 301 in release).
	svc := &mockAuthSvc{verifyEmailFn: func(_ context.Context, _ string) error { return nil }}
	r := newTestRouter(handlers.NewAuthHandler(svc, silentLogger()))

	req := httptest.NewRequest(http.MethodGet, "/auth/verify-email/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestVerifyEmail_ExpiredToken(t *testing.T) {
	// Service returns an error whose message contains "token is expired" (golang-jwt phrasing).
	svc := &mockAuthSvc{verifyEmailFn: func(_ context.Context, _ string) error {
		return fmt.Errorf("failed to parse token: token is expired")
	}}
	r := newTestRouter(handlers.NewAuthHandler(svc, silentLogger()))

	w := doRequest(r, http.MethodGet, "/auth/verify-email/expired-token", nil)

	assert.Equal(t, http.StatusGone, w.Code)
	body := parseBody(t, w)
	assert.Equal(t, "Verification token has expired", body["error"])
}

func TestVerifyEmail_InvalidToken(t *testing.T) {
	// Any error that does NOT contain "token is expired" → 400 invalid.
	svc := &mockAuthSvc{verifyEmailFn: func(_ context.Context, _ string) error {
		return fmt.Errorf("invalid or expired verification token")
	}}
	r := newTestRouter(handlers.NewAuthHandler(svc, silentLogger()))

	w := doRequest(r, http.MethodGet, "/auth/verify-email/bad-token", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	body := parseBody(t, w)
	assert.Equal(t, "Invalid verification token", body["error"])
}

// Regression: before the fix, ANY error (including "invalid or expired …")
// was caught by the "expired" check because the message contained that word.
func TestVerifyEmail_InvalidToken_NotMistakenForExpired(t *testing.T) {
	svc := &mockAuthSvc{verifyEmailFn: func(_ context.Context, _ string) error {
		return fmt.Errorf("invalid or expired email verification token")
	}}
	r := newTestRouter(handlers.NewAuthHandler(svc, silentLogger()))

	w := doRequest(r, http.MethodGet, "/auth/verify-email/some-token", nil)

	// Must be 400 (invalid), NOT 410 (expired)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ---------------------------------------------------------------------------
// ResendVerification handler tests
// ---------------------------------------------------------------------------

func TestResendVerification_Success(t *testing.T) {
	called := false
	svc := &mockAuthSvc{resendVerificationFn: func(_ context.Context, email string) error {
		called = true
		assert.Equal(t, "user@example.com", email)
		return nil
	}}
	r := newTestRouter(handlers.NewAuthHandler(svc, silentLogger()))

	w := doRequest(r, http.MethodPost, "/auth/resend-verification", map[string]string{"email": "user@example.com"})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, called)
}

func TestResendVerification_AlwaysReturns200(t *testing.T) {
	// Even when the service returns an error (unknown email, already verified, etc.)
	// the handler must return 200 to prevent email enumeration.
	svc := &mockAuthSvc{resendVerificationFn: func(_ context.Context, _ string) error {
		return fmt.Errorf("user not found")
	}}
	r := newTestRouter(handlers.NewAuthHandler(svc, silentLogger()))

	w := doRequest(r, http.MethodPost, "/auth/resend-verification", map[string]string{"email": "nobody@example.com"})

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestResendVerification_MissingEmail(t *testing.T) {
	svc := &mockAuthSvc{resendVerificationFn: func(_ context.Context, _ string) error { return nil }}
	r := newTestRouter(handlers.NewAuthHandler(svc, silentLogger()))

	w := doRequest(r, http.MethodPost, "/auth/resend-verification", map[string]string{})

	assert.Equal(t, http.StatusBadRequest, w.Code)
	body := parseBody(t, w)
	assert.Contains(t, body["error"], "email")
}

// ---------------------------------------------------------------------------
// Login handler tests
// ---------------------------------------------------------------------------

func TestLogin_EmailNotVerified_Returns403(t *testing.T) {
	svc := &mockAuthSvc{loginFn: func(_ context.Context, _ *domain.LoginRequest) (*domain.AuthResponse, error) {
		return nil, fmt.Errorf("EMAIL_NOT_VERIFIED")
	}}
	r := newTestRouter(handlers.NewAuthHandler(svc, silentLogger()))

	w := doRequest(r, http.MethodPost, "/auth/login", map[string]string{
		"email": "user@example.com", "password": "secret",
	})

	assert.Equal(t, http.StatusForbidden, w.Code)
	body := parseBody(t, w)
	assert.Equal(t, "EMAIL_NOT_VERIFIED", body["error"])
}

func TestLogin_TwoFARequired_Returns200WithFlag(t *testing.T) {
	svc := &mockAuthSvc{loginFn: func(_ context.Context, _ *domain.LoginRequest) (*domain.AuthResponse, error) {
		return nil, fmt.Errorf("2FA_REQUIRED")
	}}
	r := newTestRouter(handlers.NewAuthHandler(svc, silentLogger()))

	w := doRequest(r, http.MethodPost, "/auth/login", map[string]string{
		"email": "user@example.com", "password": "secret",
	})

	assert.Equal(t, http.StatusOK, w.Code)
	body := parseBody(t, w)
	assert.Equal(t, true, body["twofa_required"])
}

func TestLogin_InvalidCredentials_Returns401(t *testing.T) {
	svc := &mockAuthSvc{loginFn: func(_ context.Context, _ *domain.LoginRequest) (*domain.AuthResponse, error) {
		return nil, fmt.Errorf("invalid email or password")
	}}
	r := newTestRouter(handlers.NewAuthHandler(svc, silentLogger()))

	w := doRequest(r, http.MethodPost, "/auth/login", map[string]string{
		"email": "user@example.com", "password": "wrong",
	})

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogin_BadRequestBody(t *testing.T) {
	svc := &mockAuthSvc{}
	r := newTestRouter(handlers.NewAuthHandler(svc, silentLogger()))

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ---------------------------------------------------------------------------
// Register handler tests
// ---------------------------------------------------------------------------

func TestRegister_Success(t *testing.T) {
	svc := &mockAuthSvc{registerFn: func(_ context.Context, _ *domain.RegisterRequest) (*domain.AuthResponse, error) {
		return &domain.AuthResponse{
			User: domain.UserResponse{ID: "new-user-id", Email: "new@example.com"},
		}, nil
	}}
	r := newTestRouter(handlers.NewAuthHandler(svc, silentLogger()))

	w := doRequest(r, http.MethodPost, "/auth/register", map[string]string{
		"email": "new@example.com", "password": "Secret123!", "first_name": "Test", "last_name": "User",
	})

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRegister_DuplicateEmail_Returns409(t *testing.T) {
	svc := &mockAuthSvc{registerFn: func(_ context.Context, _ *domain.RegisterRequest) (*domain.AuthResponse, error) {
		return nil, fmt.Errorf("user with email already exists")
	}}
	r := newTestRouter(handlers.NewAuthHandler(svc, silentLogger()))

	w := doRequest(r, http.MethodPost, "/auth/register", map[string]string{
		"email": "existing@example.com", "password": "Secret123!", "first_name": "A", "last_name": "B",
	})

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestRegister_BadRequestBody(t *testing.T) {
	svc := &mockAuthSvc{}
	r := newTestRouter(handlers.NewAuthHandler(svc, silentLogger()))

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString("{bad json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
