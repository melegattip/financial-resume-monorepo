package domain

import "time"

// Tenant represents a multi-tenant group (family, company, team, etc.)
type Tenant struct {
	ID        string
	Name      string
	Slug      string
	OwnerID   string
	IsActive  bool
	Plan      string
	Settings  string // JSON-encoded settings
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// UpdateTenantRequest holds the fields that can be updated on a tenant.
// All fields are optional (pointer = not provided).
type UpdateTenantRequest struct {
	Name     *string `json:"name"`
	Settings *string `json:"settings"`
}
