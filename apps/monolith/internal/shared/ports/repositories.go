package ports

import "context"

// BaseRepository defines common CRUD operations.
// Module-specific repositories will embed or extend this pattern.
type BaseRepository[T any] interface {
	// FindByID retrieves an entity by its unique identifier.
	FindByID(ctx context.Context, id string) (*T, error)
	// FindAll retrieves all entities for a given user.
	FindAll(ctx context.Context, userID string) ([]T, error)
	// Create persists a new entity.
	Create(ctx context.Context, entity *T) error
	// Update modifies an existing entity.
	Update(ctx context.Context, entity *T) error
	// Delete performs a soft delete on an entity.
	Delete(ctx context.Context, id string) error
}
