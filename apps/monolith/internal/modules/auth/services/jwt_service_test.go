package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-secret-key-for-jwt-tests-2024"

var (
	testUserID1 = "00000000-0000-0000-0000-000000000001"
	testUserID2 = "00000000-0000-0000-0000-000000000042"
)

func newTestJWTService() *jwtService {
	return NewJWTService(testSecret, 1*time.Hour, 7*24*time.Hour, "test-issuer").(*jwtService)
}

func TestGenerateTokens(t *testing.T) {
	svc := newTestJWTService()

	pair, err := svc.GenerateTokens(testUserID1, "user@example.com")
	require.NoError(t, err)
	assert.NotEmpty(t, pair.AccessToken)
	assert.NotEmpty(t, pair.RefreshToken)
	assert.Equal(t, "Bearer", pair.TokenType)
	assert.True(t, pair.ExpiresAt.After(time.Now()))
}

func TestValidateAccessToken(t *testing.T) {
	svc := newTestJWTService()

	pair, err := svc.GenerateTokens(testUserID2, "test@example.com")
	require.NoError(t, err)

	claims, err := svc.ValidateAccessToken(pair.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, testUserID2, claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, "access", claims.TokenType)
}

func TestValidateRefreshToken(t *testing.T) {
	svc := newTestJWTService()

	pair, err := svc.GenerateTokens(testUserID2, "test@example.com")
	require.NoError(t, err)

	claims, err := svc.ValidateRefreshToken(pair.RefreshToken)
	require.NoError(t, err)
	assert.Equal(t, testUserID2, claims.UserID)
	assert.Equal(t, "refresh", claims.TokenType)
}

func TestValidateAccessToken_WrongType(t *testing.T) {
	svc := newTestJWTService()

	pair, err := svc.GenerateTokens(testUserID1, "user@example.com")
	require.NoError(t, err)

	// Try to validate refresh token as access token
	_, err = svc.ValidateAccessToken(pair.RefreshToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token type")
}

func TestValidateRefreshToken_WrongType(t *testing.T) {
	svc := newTestJWTService()

	pair, err := svc.GenerateTokens(testUserID1, "user@example.com")
	require.NoError(t, err)

	// Try to validate access token as refresh token
	_, err = svc.ValidateRefreshToken(pair.AccessToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token type")
}

func TestValidateAccessToken_InvalidSignature(t *testing.T) {
	svc := newTestJWTService()
	otherSvc := NewJWTService("different-secret", 1*time.Hour, 7*24*time.Hour, "test-issuer")

	pair, err := otherSvc.GenerateTokens(testUserID1, "user@example.com")
	require.NoError(t, err)

	_, err = svc.ValidateAccessToken(pair.AccessToken)
	assert.Error(t, err)
}

func TestValidateAccessToken_Expired(t *testing.T) {
	// Create service with very short expiry
	svc := NewJWTService(testSecret, -1*time.Second, 7*24*time.Hour, "test-issuer")

	pair, err := svc.GenerateTokens(testUserID1, "user@example.com")
	require.NoError(t, err)

	_, err = svc.ValidateAccessToken(pair.AccessToken)
	assert.Error(t, err)
}

func TestValidateAccessToken_Malformed(t *testing.T) {
	svc := newTestJWTService()

	_, err := svc.ValidateAccessToken("not-a-valid-jwt")
	assert.Error(t, err)
}

func TestGenerateEmailVerificationToken(t *testing.T) {
	svc := newTestJWTService()

	token, err := svc.GenerateEmailVerificationToken(testUserID1, "user@example.com")
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := svc.ValidateEmailVerificationToken(token)
	require.NoError(t, err)
	assert.Equal(t, testUserID1, claims.UserID)
	assert.Equal(t, "email_verification", claims.TokenType)
}

func TestGeneratePasswordResetToken(t *testing.T) {
	svc := newTestJWTService()

	token, err := svc.GeneratePasswordResetToken(testUserID1, "user@example.com")
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := svc.ValidatePasswordResetToken(token)
	require.NoError(t, err)
	assert.Equal(t, testUserID1, claims.UserID)
	assert.Equal(t, "password_reset", claims.TokenType)
}

func TestValidateEmailVerificationToken_WrongType(t *testing.T) {
	svc := newTestJWTService()

	token, err := svc.GeneratePasswordResetToken(testUserID1, "user@example.com")
	require.NoError(t, err)

	_, err = svc.ValidateEmailVerificationToken(token)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token type")
}

func TestClaimsContainIssuer(t *testing.T) {
	svc := newTestJWTService()

	pair, err := svc.GenerateTokens(testUserID1, "user@example.com")
	require.NoError(t, err)

	claims, err := svc.ValidateAccessToken(pair.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, "test-issuer", claims.Issuer)
}
