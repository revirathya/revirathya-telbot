package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	financecorev1 "github.com/revirathya/revirathya-telbot/gen/go/financecore/v1"
	"github.com/revirathya/revirathya-telbot/internal/domain"
	grpcCon "github.com/revirathya/revirathya-telbot/internal/platform/grpc"
)

type expenseRepositoryImpl struct {
	fincoreConn *grpcCon.FinanceCoreConn
}

func NewExpenseRepositoryImpl(fincoreConn *grpcCon.FinanceCoreConn) domain.ExpenseRepository {
	return &expenseRepositoryImpl{
		fincoreConn: fincoreConn,
	}
}

func (r *expenseRepositoryImpl) CreateExpense(ctx context.Context, newExpense domain.CreateExpense) (*domain.Expense, error) {
	slog.Info("test from repo", "data", newExpense)
	in := &financecorev1.CreateExpenseRequest{
		UserId:        newExpense.UserID.String(),
		CategoryId:    newExpense.CategoryID.String(),
		TransactionDt: newExpense.TransactionDate.Format("2006-01-02"),
		Description:   newExpense.Description,
		Amount:        strconv.FormatInt(newExpense.Amount, 10),
	}

	row, err := r.fincoreConn.Expense.CreateExpense(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed to create expenses: %v", err)
	}

	uuidID, err := uuid.Parse(row.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid id UUID format: %v", err)
	}

	uuidUserID, err := uuid.Parse(row.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid user id UUID format: %v", err)
	}

	dtTransactionDt, err := time.Parse("2006-01-02", row.TransactionDt)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction date Date format: %v", err)
	}

	uuidCategoryID, err := uuid.Parse(row.CategoryId)
	if err != nil {
		return nil, fmt.Errorf("invalid category id UUID format: %v", err)
	}

	decAmount, err := decimal.NewFromString(row.Amount)
	if err != nil {
		return nil, fmt.Errorf("invalid amount Decimal format: %v", err)
	}

	result := &domain.Expense{
		ID:              uuidID,
		UserID:          uuidUserID,
		TransactionDate: dtTransactionDt,
		CategoryID:      uuidCategoryID,
		Description:     row.Description,
		Amount:          decAmount,
	}
	return result, nil
}

func (r *expenseRepositoryImpl) ListExpenses(ctx context.Context, userID uuid.UUID, transactionDate time.Time) ([]domain.Expense, error) {
	trxDate := transactionDate.Format("2006-01-02")
	in := &financecorev1.ListExpensesRequest{
		UserId:        userID.String(),
		TransactionDt: &trxDate,
	}

	rows, err := r.fincoreConn.Expense.ListExpenses(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed to list expenses: %v", err)
	}

	results := make([]domain.Expense, 0, len(rows.Expenses))
	for _, row := range rows.Expenses {
		uuidID, err := uuid.Parse(row.Id)
		if err != nil {
			return nil, fmt.Errorf("invalid id UUID format: %v", err)
		}

		uuidUserID, err := uuid.Parse(row.UserId)
		if err != nil {
			return nil, fmt.Errorf("invalid user id UUID format: %v", err)
		}

		dtTransactionDt, err := time.Parse("2006-01-02", row.TransactionDt)
		if err != nil {
			return nil, fmt.Errorf("invalid transaction date Date format: %v", err)
		}

		uuidCategoryID, err := uuid.Parse(row.CategoryId)
		if err != nil {
			return nil, fmt.Errorf("invalid category id UUID format: %v", err)
		}

		decAmount, err := decimal.NewFromString(row.Amount)
		if err != nil {
			return nil, fmt.Errorf("invalid amount Decimal format: %v", err)
		}

		results = append(results, domain.Expense{
			ID:              uuidID,
			UserID:          uuidUserID,
			TransactionDate: dtTransactionDt,
			CategoryID:      uuidCategoryID,
			CategoryName:    row.CategoryName,
			Description:     row.Description,
			Amount:          decAmount,
		})
	}
	return results, nil
}

func (r *expenseRepositoryImpl) GetExpenseByID(ctx context.Context, id uuid.UUID) (*domain.Expense, error) {
	in := &financecorev1.GetExpenseByIDRequest{
		Id: id.String(),
	}

	row, err := r.fincoreConn.Expense.GetExpenseByID(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed to get expense by id: %v", err)
	}

	uuidID, err := uuid.Parse(row.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid id UUID format: %v", err)
	}

	uuidUserID, err := uuid.Parse(row.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid user id UUID format: %v", err)
	}

	dtTransactionDt, err := time.Parse("2006-01-02", row.TransactionDt)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction date Date format: %v", err)
	}

	uuidCategoryID, err := uuid.Parse(row.CategoryId)
	if err != nil {
		return nil, fmt.Errorf("invalid category id UUID format: %v", err)
	}

	decAmount, err := decimal.NewFromString(row.Amount)
	if err != nil {
		return nil, fmt.Errorf("invalid amount Decimal format: %v", err)
	}

	result := &domain.Expense{
		ID:              uuidID,
		UserID:          uuidUserID,
		TransactionDate: dtTransactionDt,
		CategoryID:      uuidCategoryID,
		Description:     row.Description,
		Amount:          decAmount,
	}
	return result, nil
}

func (r *expenseRepositoryImpl) DeleteExpense(ctx context.Context, id uuid.UUID) error {
	in := &financecorev1.DeleteExpenseRequest{
		Id: id.String(),
	}

	if _, err := r.fincoreConn.Expense.DeleteExpense(ctx, in); err != nil {
		return fmt.Errorf("failed to delete expense: %v", err)
	}
	return nil
}
