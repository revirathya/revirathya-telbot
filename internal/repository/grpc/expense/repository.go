package expense

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	financecorev1 "github.com/revirathya/revirathya-finance-core/gen/grpc/financecore/v1"

	// Internal
	"github.com/revirathya/revirathya-telbot/internal/domain"
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

func (r *expenseRepositoryImpl) CreateExpense(ctx context.Context, form domain.ExpenseForm, requestID uuid.UUID) (*domain.Expense, error) {
	in := &financecorev1.CreateExpenseRequest{
		CategoryId:    form.CategoryID.String(),
		TransactionDt: converter.ToProtoDate(form.TransactionDt),
		Description:   &form.Description,
		Amount:        converter.Int64ToProtoDecimal(form.Amount),
	}

	row, err := r.finance.Expense.CreateExpense(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed to create expenses: %v", err)
	}

	id, err := converter.ParseUUID("id", row.Expense.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid id UUID format: %v", err)
	}

	transactionDt, err := converter.FromProtoDate("transaction_dt", row.Expense.TransactionDt)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction date Date format: %v", err)
	}

	categoryID, err := converter.ParseUUID("category_id", row.Expense.CategoryId)
	if err != nil {
		return nil, fmt.Errorf("invalid category id UUID format: %v", err)
	}

	amount, err := converter.FromProtoDecimal("amount", row.Expense.Amount)
	if err != nil {
		return nil, fmt.Errorf("invalid amount Decimal format: %v", err)
	}

	result := &domain.Expense{
		ID:            id,
		TransactionDt: transactionDt,
		CategoryID:    categoryID,
		Description:   row.Expense.GetDescription(),
		Amount:        amount,
	}
	return result, nil
}

func (r *expenseRepositoryImpl) ListExpenses(ctx context.Context, transactionDt time.Time) ([]domain.Expense, error) {
	in := &financecorev1.ListExpensesRequest{
		TransactionDt: converter.ToProtoDate(transactionDt),
	}

	rows, err := r.finance.Expense.ListExpenses(ctx, in)
	if err != nil {
		return nil, err
	}

	results := make([]domain.Expense, 0, len(rows.Expenses))
	for _, row := range rows.Expenses {
		id, err := converter.ParseUUID("id", row.Id)
		if err != nil {
			return nil, err
		}

		transactionDt, err := converter.FromProtoDate("transaction_dt", row.TransactionDt)
		if err != nil {
			return nil, err
		}

		categoryID, err := converter.ParseUUID("category_id", row.CategoryId)
		if err != nil {
			return nil, fmt.Errorf("invalid category id UUID format: %v", err)
		}

		amount, err := converter.FromProtoDecimal("amount", row.Amount)
		if err != nil {
			return nil, fmt.Errorf("invalid amount Decimal format: %v", err)
		}

		results = append(results, domain.Expense{
			ID:            id,
			TransactionDt: transactionDt,
			CategoryID:    categoryID,
			CategoryName:  row.CategoryName,
			Description:   row.Description,
			Amount:        amount,
		})
	}
	return results, nil
}

func (r *expenseRepositoryImpl) GetExpenseByID(ctx context.Context, id uuid.UUID) (*domain.Expense, error) {
	in := &financecorev1.GetExpenseByIDRequest{
		Id: id.String(),
	}

	row, err := r.finance.Expense.GetExpenseByID(ctx, in)
	if err != nil {
		return nil, err
	}

	transactionDt, err := converter.FromProtoDate("transaction_dt", row.Expense.TransactionDt)
	if err != nil {
		return nil, err
	}

	categoryID, err := converter.ParseUUID("category_id", row.Expense.CategoryId)
	if err != nil {
		return nil, err
	}

	amount, err := converter.FromProtoDecimal("amount", row.Expense.Amount)
	if err != nil {
		return nil, err
	}

	result := &domain.Expense{
		ID:            id,
		TransactionDt: transactionDt,
		CategoryID:    categoryID,
		CategoryName:  row.Expense.CategoryName,
		Description:   row.Expense.Description,
		Amount:        amount,
	}
	return result, nil
}

func (r *expenseRepositoryImpl) DeleteExpense(ctx context.Context, id uuid.UUID) error {
	in := &financecorev1.DeleteExpenseRequest{
		Id: id.String(),
	}

	if _, err := r.finance.Expense.DeleteExpense(ctx, in); err != nil {
		return fmt.Errorf("failed to delete expense: %v", err)
	}
	return nil
}
