package category

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	financecorev1 "github.com/revirathya/revirathya-finance-core/gen/grpc/financecore/v1"

	// Internal
	"github.com/revirathya/revirathya-telbot/internal/domain"
	"github.com/revirathya/revirathya-telbot/internal/platform/grpc/clients"
)

type categoryRepositoryImpl struct {
	finance *clients.FinanceClient
}

func NewCategoryRepositoryImpl(finance *clients.FinanceClient) domain.CategoryRepository {
	return &categoryRepositoryImpl{
		finance: finance,
	}
}

func (r *categoryRepositoryImpl) ListCategories(ctx context.Context) ([]domain.Category, error) {
	in := &financecorev1.ListCategoriesRequest{}
	rows, err := r.finance.Category.ListCategories(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %v", err)
	}

	results := make([]domain.Category, 0, len(rows.Categories))
	for _, row := range rows.Categories {
		id, err := uuid.Parse(row.Id)
		if err != nil {
			return nil, fmt.Errorf("invalid id UUID format: %v", err)
		}

		results = append(results, domain.Category{
			ID:           id,
			CategoryName: row.CategoryName,
		})
	}
	return results, nil
}

func (r *categoryRepositoryImpl) GetCategoryByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	in := &financecorev1.GetCategoryByIDRequest{
		Id: id.String(),
	}

	row, err := r.finance.Category.GetCategoryByID(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed to get category by id: %v", err)
	}

	result := &domain.Category{
		ID:           id,
		CategoryName: row.Category.CategoryName,
	}
	return result, nil
}
