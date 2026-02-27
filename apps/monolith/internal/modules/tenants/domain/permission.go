package domain

// Permission represents a single permission entry in the global catalog.
type Permission struct {
	Key         string
	Description string
	Category    string
}

// RolePermission links a role to a permission within a specific tenant.
type RolePermission struct {
	ID            string
	TenantID      string
	Role          string
	PermissionKey string
}
