package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsUsable_ActiveNoLimits(t *testing.T) {
	inv := &Invitation{
		IsActive: true,
		MaxUses:  0, // 0 = unlimited
	}
	assert.True(t, inv.IsUsable())
}

func TestIsUsable_Inactive(t *testing.T) {
	inv := &Invitation{IsActive: false}
	assert.False(t, inv.IsUsable())
}

func TestIsUsable_MaxUsesReached(t *testing.T) {
	inv := &Invitation{
		IsActive:  true,
		MaxUses:   5,
		UsedCount: 5,
	}
	assert.False(t, inv.IsUsable())
}

func TestIsUsable_MaxUsesNotReached(t *testing.T) {
	inv := &Invitation{
		IsActive:  true,
		MaxUses:   5,
		UsedCount: 4,
	}
	assert.True(t, inv.IsUsable())
}

func TestIsUsable_Expired(t *testing.T) {
	past := time.Now().Add(-1 * time.Hour)
	inv := &Invitation{
		IsActive:  true,
		ExpiresAt: &past,
	}
	assert.False(t, inv.IsUsable())
}

func TestIsUsable_NotExpired(t *testing.T) {
	future := time.Now().Add(24 * time.Hour)
	inv := &Invitation{
		IsActive:  true,
		ExpiresAt: &future,
	}
	assert.True(t, inv.IsUsable())
}

func TestIsUsable_NoExpiryDate(t *testing.T) {
	inv := &Invitation{
		IsActive:  true,
		ExpiresAt: nil,
	}
	assert.True(t, inv.IsUsable())
}

func TestGenerateInviteCode_Length(t *testing.T) {
	code := GenerateInviteCode()
	assert.Len(t, code, 8)
}

func TestGenerateInviteCode_AlphanumericOnly(t *testing.T) {
	const validChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i := 0; i < 20; i++ {
		code := GenerateInviteCode()
		for _, ch := range code {
			assert.Contains(t, validChars, string(ch), "code contains invalid character: %c", ch)
		}
	}
}

func TestGenerateInviteCode_ReasonablyUnique(t *testing.T) {
	codes := make(map[string]bool)
	for i := 0; i < 100; i++ {
		codes[GenerateInviteCode()] = true
	}
	assert.Greater(t, len(codes), 50, "codes should be reasonably unique across 100 generations")
}
