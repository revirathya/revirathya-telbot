package grpc

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	financecorev1 "github.com/revirathya/revirathya-telbot/gen/go/financecore/v1"
	"github.com/revirathya/revirathya-telbot/internal/domain"
	grpcCon "github.com/revirathya/revirathya-telbot/internal/platform/grpc"
)

type categoryRepositoryImpl struct {
	fincoreConn *grpcCon.FinanceCoreConn
}

func NewCategoryRepositoryImpl(fincoreConn *grpcCon.FinanceCoreConn) domain.CategoryRepository {
	return &categoryRepositoryImpl{
		fincoreConn: fincoreConn,
	}
}

func (r *categoryRepositoryImpl) ListCategories(ctx context.Context, userID uuid.UUID) ([]domain.Category, error) {
	in := &financecorev1.ListCategoriesRequest{}
	rows, err := r.fincoreConn.Category.ListCategories(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %v", err)
	}

	results := make([]domain.Category, 0, len(rows.Categories))
	for _, row := range rows.Categories {
		uuidID, err := uuid.Parse(row.Id)
		if err != nil {
			return nil, fmt.Errorf("invalid id UUID format: %v", err)
		}

		results = append(results, domain.Category{
			ID:           uuidID,
			CategoryName: row.CategoryName,
		})
	}
	return results, nil
}

func (r *categoryRepositoryImpl) GetCategoryByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	in := &financecorev1.GetCategoryByIDRequest{
		Id: id.String(),
	}

	row, err := r.fincoreConn.Category.GetCategoryByID(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed to get category by id: %v", err)
	}

	uuidID, err := uuid.Parse(row.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid id UUID format: %v", err)
	}

	category := &domain.Category{
		ID:           uuidID,
		CategoryName: row.CategoryName,
	}
	return category, nil
}
