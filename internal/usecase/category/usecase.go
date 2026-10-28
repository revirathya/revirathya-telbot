package category

import (
	"context"

	"github.com/google/uuid"
	"github.com/revirathya/revirathya-telbot/internal/domain"
)

// Expense: Usecase
type CategoryUsecase struct {
	category domain.CategoryRepository
}

func NewCategoryUsecase(category domain.CategoryRepository) *CategoryUsecase {
	return &CategoryUsecase{
		category: category,
	}
}

func (uc *CategoryUsecase) ListCategories(ctx context.Context) ([]domain.Category, error) {
	return uc.category.ListCategories(ctx)
}

func (uc *CategoryUsecase) GetCategoryByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	return uc.category.GetCategoryByID(ctx, id)
}
