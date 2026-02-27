package domain

import (
	"math/rand"
	"time"
)

// Invitation represents a tenant invitation link (code-based).
type Invitation struct {
	ID        string
	TenantID  string
	Code      string
	Role      string
	CreatedBy string
	ExpiresAt *time.Time
	MaxUses   int
	UsedCount int
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// IsUsable returns true if the invitation can still be accepted.
func (inv *Invitation) IsUsable() bool {
	if !inv.IsActive {
		return false
	}
	if inv.MaxUses > 0 && inv.UsedCount >= inv.MaxUses {
		return false
	}
	if inv.ExpiresAt != nil && time.Now().After(*inv.ExpiresAt) {
		return false
	}
	return true
}

// GenerateInviteCode creates a random 8-character alphanumeric invitation code.
func GenerateInviteCode() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

// CreateInvitationRequest holds the data to create a new invitation.
type CreateInvitationRequest struct {
	Role      string     `json:"role" binding:"required,oneof=admin member viewer"`
	ExpiresAt *time.Time `json:"expires_at"`
	MaxUses   int        `json:"max_uses"`
}

// JoinTenantRequest holds the invitation code used to join a tenant.
type JoinTenantRequest struct {
	Code string `json:"code" binding:"required"`
}
