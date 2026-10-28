package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	// Internal
	"github.com/revirathya/revirathya-telbot/internal/domain/validation"
)

// Entity
type Expense struct {
	ID            uuid.UUID
	TransactionDt time.Time
	CategoryID    uuid.UUID
	CategoryName  string
	Description   string
	Amount        decimal.Decimal
}

// Form
type ExpenseForm struct {
	CategoryID    uuid.UUID
	TransactionDt time.Time
	Description   string
	Amount        int64
}

func (f *ExpenseForm) Validate() error {
	// Amount
	if err := validation.ValidateInt64("amount", f.Amount, true); err != nil {
		return err
	}

	// Date
	if err := validation.ValidateDate("transaction_dt", f.TransactionDt); err != nil {
		return err
	}

	// Description
	description, err := validation.ValidateString("description", f.Description, 500, false)
	if err != nil {
		return err
	}
	f.Description = description

	return nil
}

// Repository: Expense
type ExpenseRepository interface {
	CreateExpense(ctx context.Context, form ExpenseForm, requestID uuid.UUID) (*Expense, error)
	ListExpenses(ctx context.Context, transactionDate time.Time) ([]Expense, error)
	GetExpenseByID(ctx context.Context, id uuid.UUID) (*Expense, error)
	DeleteExpense(ctx context.Context, id uuid.UUID) error
}
