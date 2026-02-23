package repository

import (
	"context"

	"github.com/google/uuid"
)

// MaterialRepository defines persistence operations for Material
type MaterialRepository interface {
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
