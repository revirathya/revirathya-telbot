package expense

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/revirathya/revirathya-telbot/internal/domain"
)

// Usecase: Expense
type ExpenseUsecase struct {
	expense domain.ExpenseRepository
}

func NewExpenseUsecase(expense domain.ExpenseRepository) *ExpenseUsecase {
	return &ExpenseUsecase{
		expense: expense,
	}
}

func (uc *ExpenseUsecase) CreateExpense(ctx context.Context, form domain.ExpenseForm, requestID uuid.UUID) (*domain.Expense, error) {
	return uc.expense.CreateExpense(ctx, form, requestID)
}

func (uc *ExpenseUsecase) ListExpenses(ctx context.Context, transactionDate time.Time) ([]domain.Expense, error) {
	return uc.expense.ListExpenses(ctx, transactionDate)
}

func (uc *ExpenseUsecase) GetExpenseByID(ctx context.Context, id uuid.UUID) (*domain.Expense, error) {
	return uc.expense.GetExpenseByID(ctx, id)
}

func (uc *ExpenseUsecase) DeleteExpense(ctx context.Context, id uuid.UUID) error {
	return uc.expense.DeleteExpense(ctx, id)
}
