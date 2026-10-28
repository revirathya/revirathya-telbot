package expense

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	financecorev1 "github.com/revirathya/revirathya-finance-core/gen/grpc/financecore/v1"

	// Internal
	"github.com/revirathya/revirathya-telbot/internal/domain"
	errDom "github.com/revirathya/revirathya-telbot/internal/domain/error"
	"github.com/revirathya/revirathya-telbot/internal/platform/grpc/clients"
	"github.com/revirathya/revirathya-telbot/internal/platform/grpc/converter"
)

type expenseRepositoryImpl struct {
	finance *clients.FinanceClient
}

func NewExpenseRepositoryImpl(finance *clients.FinanceClient) domain.ExpenseRepository {
	return &expenseRepositoryImpl{
		finance: finance,
	}
}

func (r *expenseRepositoryImpl) expenseToDomain(exp *financecorev1.Expense) (*domain.Expense, error) {
	id, err := converter.FromProtoUUID("id", exp.Id)
	if err != nil {
		return nil, err
	}

	transactionDt, err := converter.FromProtoDate("transaction_dt", exp.TransactionDt)
	if err != nil {
		return nil, err
	}

	categoryID, err := converter.FromProtoUUID("category_id", exp.CategoryId)
	if err != nil {
		return nil, err
	}

	amount, err := converter.FromProtoDecimal("amount", exp.Amount)
	if err != nil {
		return nil, err
	}

	result := &domain.Expense{
		ID:            id,
		TransactionDt: transactionDt,
		CategoryID:    categoryID,
		Description:   exp.GetDescription(),
		Amount:        amount,
	}
	return result, nil
}

// Public
func (r *expenseRepositoryImpl) CreateExpense(ctx context.Context, form domain.ExpenseForm, requestID uuid.UUID) (*domain.Expense, error) {
	// Prepare input
	in := &financecorev1.CreateExpenseRequest{
		CategoryId:    form.CategoryID.String(),
		TransactionDt: converter.ToProtoDate(form.TransactionDt),
		Description:   &form.Description,
		Amount:        converter.Int64ToProtoDecimal(form.Amount),
	}

	// Create expense
	row, err := r.finance.Expense.CreateExpense(ctx, in)
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			return nil, errDom.ErrInvalidInput
		}
		return nil, fmt.Errorf("failed to create expenses: %v", err)
	}

	// Parse response
	return r.expenseToDomain(row.Expense)
}

func (r *expenseRepositoryImpl) ListExpenses(ctx context.Context, transactionDt time.Time) ([]domain.Expense, error) {
	// Prepare input
	in := &financecorev1.ListExpensesRequest{
		TransactionDt: converter.ToProtoDate(transactionDt),
	}

	// List expenses
	rows, err := r.finance.Expense.ListExpenses(ctx, in)
	if err != nil {
		return nil, err
	}

	// Parse response
	results := make([]domain.Expense, 0, len(rows.Expenses))
	for _, row := range rows.Expenses {
		result, err := r.expenseToDomain(row)
		if err != nil {
			return nil, err
		}

		results = append(results, *result)
	}
	return results, nil
}

func (r *expenseRepositoryImpl) GetExpenseByID(ctx context.Context, id uuid.UUID) (*domain.Expense, error) {
	// Prepare input
	in := &financecorev1.GetExpenseByIDRequest{
		Id: id.String(),
	}

	// Get expense by ID
	row, err := r.finance.Expense.GetExpenseByID(ctx, in)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, errDom.ErrNotFound
		}
		return nil, err
	}

	// Parse response
	return r.expenseToDomain(row.Expense)
}

func (r *expenseRepositoryImpl) DeleteExpense(ctx context.Context, id uuid.UUID) error {
	in := &financecorev1.DeleteExpenseRequest{
		Id: id.String(),
	}

	if _, err := r.finance.Expense.DeleteExpense(ctx, in); err != nil {
		if status.Code(err) == codes.NotFound {
			return errDom.ErrNotFound
		}
		return err
	}
	return nil
}
