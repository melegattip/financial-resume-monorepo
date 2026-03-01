package domain

import "time"

// AuditLog represents an immutable record of an important action within a tenant.
type AuditLog struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	UserID      string    `json:"user_id"`
	UserName    string    `json:"user_name,omitempty"`    // populated via JOIN on read
	Action      string    `json:"action"`
	Description string    `json:"description,omitempty"` // human-readable summary of the action
	EntityType  string    `json:"entity_type"`
	EntityID    string    `json:"entity_id"`
	OldValues   string    `json:"old_values,omitempty"`
	NewValues   string    `json:"new_values,omitempty"`
	IPAddress   string    `json:"ip_address,omitempty"`
	UserAgent   string    `json:"user_agent,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
