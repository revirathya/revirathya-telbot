package domain

import (
	"context"

	"github.com/google/uuid"
)

// Entity
type Category struct {
	ID           uuid.UUID
	CategoryName string
}

// Entity
type CategoryRepository interface {
	ListCategories(ctx context.Context) ([]Category, error)
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*Category, error)
}
