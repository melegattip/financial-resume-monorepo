package domain

// API Contract tests for AuthResponse.
//
// Purpose: pin the JSON shape of auth responses so that any change to the
// response structure is caught immediately, before it reaches production.
//
// Root cause of the frontend/backend discrepancy this prevents:
//   - Old microservice returned flat: { "access_token": "...", "user": {...} }
//   - Monolith returns nested:        { "user": {...}, "tokens": { "access_token": "..." } }
//   - Without these tests, the mismatch is only found at runtime in production.

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildTestAuthResponse() AuthResponse {
	return AuthResponse{
		User: UserResponse{
			ID:        "user-123",
			Email:     "test@example.com",
			FirstName: "Test",
			LastName:  "User",
			IsActive:  true,
			IsVerified: false,
			CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		Tokens: TokenPair{
			AccessToken:  "access.token.here",
			RefreshToken: "refresh.token.here",
			ExpiresAt:    time.Date(2026, 2, 21, 0, 0, 0, 0, time.UTC),
			TokenType:    "Bearer",
		},
	}
}

// TestAuthResponse_TopLevelShape verifies the top-level keys of the JSON response.
// Regression: access_token must NOT appear at top level (it lives under "tokens").
func TestAuthResponse_TopLevelShape(t *testing.T) {
	b, err := json.Marshal(buildTestAuthResponse())
	require.NoError(t, err)

	var raw map[string]any
	require.NoError(t, json.Unmarshal(b, &raw))

	assert.Contains(t, raw, "user", "response must have top-level 'user'")
	assert.Contains(t, raw, "tokens", "response must have top-level 'tokens' (nested, not flat)")

	// Regression guard: these keys must NOT appear at the top level
	assert.NotContains(t, raw, "access_token", "access_token must be nested under tokens{}, not at top level")
	assert.NotContains(t, raw, "refresh_token", "refresh_token must be nested under tokens{}, not at top level")
	assert.NotContains(t, raw, "expires_at", "expires_at must be nested under tokens{}, not at top level")
}

// TestAuthResponse_TokensShape verifies the structure of the nested tokens object.
func TestAuthResponse_TokensShape(t *testing.T) {
	b, err := json.Marshal(buildTestAuthResponse())
	require.NoError(t, err)

	var raw map[string]any
	require.NoError(t, json.Unmarshal(b, &raw))

	tokens, ok := raw["tokens"].(map[string]any)
	require.True(t, ok, "tokens must be a JSON object")

	assert.Contains(t, tokens, "access_token")
	assert.Contains(t, tokens, "refresh_token")
	assert.Contains(t, tokens, "expires_at")
	assert.Contains(t, tokens, "token_type")

	assert.Equal(t, "access.token.here", tokens["access_token"])
	assert.Equal(t, "Bearer", tokens["token_type"])
}

// TestAuthResponse_ExpiresAt_IsISO8601String verifies that expires_at serializes as
// an ISO8601 string (time.Time), NOT a Unix integer.
//
// Frontend impact: clients must convert this string to a Unix timestamp via
// new Date(expires_at).getTime() / 1000 before calling parseInt().
func TestAuthResponse_ExpiresAt_IsISO8601String(t *testing.T) {
	b, err := json.Marshal(buildTestAuthResponse())
	require.NoError(t, err)

	var raw map[string]any
	require.NoError(t, json.Unmarshal(b, &raw))

	tokens := raw["tokens"].(map[string]any)
	expiresAt, ok := tokens["expires_at"].(string)
	require.True(t, ok, "expires_at must serialize as a string, not a number")

	_, parseErr := time.Parse(time.RFC3339, expiresAt)
	assert.NoError(t, parseErr, "expires_at must be valid RFC3339/ISO8601 format")
}

// TestAuthResponse_UserShape verifies the structure of the user object.
func TestAuthResponse_UserShape(t *testing.T) {
	b, err := json.Marshal(buildTestAuthResponse())
	require.NoError(t, err)

	var raw map[string]any
	require.NoError(t, json.Unmarshal(b, &raw))

	user, ok := raw["user"].(map[string]any)
	require.True(t, ok, "user must be a JSON object")

	assert.Contains(t, user, "id")
	assert.Contains(t, user, "email")
	assert.Contains(t, user, "first_name")
	assert.Contains(t, user, "last_name")
	assert.Contains(t, user, "is_active")
	assert.Contains(t, user, "is_verified")
	assert.Contains(t, user, "created_at")
}

// TestAuthResponse_NoSensitiveFieldsLeaked ensures passwords and internal tokens
// are never serialized in the auth response.
func TestAuthResponse_NoSensitiveFieldsLeaked(t *testing.T) {
	b, err := json.Marshal(buildTestAuthResponse())
	require.NoError(t, err)

	rawJSON := string(b)
	assert.False(t, strings.Contains(rawJSON, "password"), "password must never appear in auth response JSON")
}
