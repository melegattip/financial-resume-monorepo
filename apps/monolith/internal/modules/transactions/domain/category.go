package domain

import (
	"errors"
	"time"
)

// Category represents an expense category
type Category struct {
	ID        string
	UserID    string
	TenantID  string
	Name      string
	Color     string
	Icon      string
	Priority  int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// NewCategory creates a new category with validation
func NewCategory(userID, name, color, icon string, priority int) (*Category, error) {
	if userID == "" {
		return nil, errors.New("user_id is required")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}

	now := time.Now().UTC()
	return &Category{
		UserID:    userID,
		Name:      name,
		Color:     color,
		Icon:      icon,
		Priority:  priority,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Update updates mutable fields
func (c *Category) Update(name, color, icon string, priority int) error {
	if name == "" {
		return errors.New("name is required")
	}

	c.Name = name
	c.Color = color
	c.Icon = icon
	c.Priority = priority
	c.UpdatedAt = time.Now().UTC()
	return nil
}

// SoftDelete marks the category as deleted
func (c *Category) SoftDelete() {
	now := time.Now().UTC()
	c.DeletedAt = &now
}

// IsDeleted checks if the category has been soft-deleted
func (c *Category) IsDeleted() bool {
	return c.DeletedAt != nil
}
