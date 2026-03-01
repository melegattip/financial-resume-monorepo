package domain

import "time"

// Tenant represents a multi-tenant group (family, company, team, etc.)
type Tenant struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
	OwnerID   string     `json:"owner_id"`
	IsActive  bool       `json:"is_active"`
	Plan      string     `json:"plan"`
	Settings  string     `json:"settings"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// UpdateTenantRequest holds the fields that can be updated on a tenant.
// All fields are optional (pointer = not provided).
type UpdateTenantRequest struct {
	Name     *string `json:"name"`
	Settings *string `json:"settings"`
}
