package expense

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/revirathya/revirathya-telbot/internal/domain"
)

// Usecase: Expense
type ExpenseUsecase struct {
	expRepo domain.ExpenseRepository
}

func NewExpenseUsecase(expRepo domain.ExpenseRepository) *ExpenseUsecase {
	return &ExpenseUsecase{
		expRepo: expRepo,
	}
}

func (uc *ExpenseUsecase) CreateExpense(ctx context.Context, newExpense domain.CreateExpense) (*domain.Expense, error) {
	return uc.expRepo.CreateExpense(ctx, newExpense)
}

func (uc *ExpenseUsecase) ListExpenses(ctx context.Context, userID uuid.UUID, transactionDate time.Time) ([]domain.Expense, error) {
	return uc.expRepo.ListExpenses(ctx, userID, transactionDate)
}

func (uc *ExpenseUsecase) GetExpenseByID(ctx context.Context, id uuid.UUID) (*domain.Expense, error) {
	return uc.expRepo.GetExpenseByID(ctx, id)
}

func (uc *ExpenseUsecase) DeleteExpense(ctx context.Context, id uuid.UUID) error {
	return uc.expRepo.DeleteExpense(ctx, id)
}
