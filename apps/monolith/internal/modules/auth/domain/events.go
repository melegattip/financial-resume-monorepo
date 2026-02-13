package domain

import (
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/shared/events"
)

const (
	EventUserRegistered      = "user.registered"
	EventUserLoggedIn        = "user.logged_in"
	EventUserPasswordChanged = "user.password_changed"
	EventUser2FAEnabled      = "user.2fa_enabled"
	EventUserDeleted         = "user.deleted"
)

// NewUserRegisteredEvent creates a domain event for user registration.
func NewUserRegisteredEvent(userID string, email string) events.DomainEvent {
	return events.NewDomainEvent(
		EventUserRegistered,
		userID,
		userID,
		map[string]string{"email": email},
	)
}

// NewUserLoggedInEvent creates a domain event for user login.
func NewUserLoggedInEvent(userID string) events.DomainEvent {
	return events.NewDomainEvent(
		EventUserLoggedIn,
		userID,
		userID,
		nil,
	)
}

// NewUserPasswordChangedEvent creates a domain event for password change.
func NewUserPasswordChangedEvent(userID string) events.DomainEvent {
	return events.NewDomainEvent(
		EventUserPasswordChanged,
		userID,
		userID,
		nil,
	)
}

// NewUser2FAEnabledEvent creates a domain event for 2FA activation.
func NewUser2FAEnabledEvent(userID string) events.DomainEvent {
	return events.NewDomainEvent(
		EventUser2FAEnabled,
		userID,
		userID,
		nil,
	)
}

// NewUserDeletedEvent creates a domain event for account deletion.
func NewUserDeletedEvent(userID string, email string) events.DomainEvent {
	return events.NewDomainEvent(
		EventUserDeleted,
		userID,
		userID,
		map[string]string{"email": email},
	)
}
