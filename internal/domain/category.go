package domain

import (
	"context"

	"github.com/google/uuid"
)

// Entity: Category
type Category struct {
	ID           uuid.UUID
	CategoryName string
}

// Entity: Category
type CategoryRepository interface {
	ListCategories(ctx context.Context, userID uuid.UUID) ([]Category, error)
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*Category, error)
}
