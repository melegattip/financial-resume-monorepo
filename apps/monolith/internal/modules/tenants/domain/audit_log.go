package domain

import "time"

// AuditLog represents an immutable record of an important action within a tenant.
type AuditLog struct {
	ID         string
	TenantID   string
	UserID     string
	Action     string
	EntityType string
	EntityID   string
	OldValues  string // JSON-encoded
	NewValues  string // JSON-encoded
	IPAddress  string
	UserAgent  string
	CreatedAt  time.Time
}
