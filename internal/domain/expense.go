package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Entity: Expense
type Expense struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	TransactionDate time.Time
	CategoryID      uuid.UUID
	CategoryName    string
	Description     *string
	Amount          decimal.Decimal
}

type CreateExpense struct {
	UserID          uuid.UUID
	TransactionDate time.Time
	CategoryID      uuid.UUID
	Description     *string
	Amount          int64
}

// Repository: Expense
type ExpenseRepository interface {
	CreateExpense(ctx context.Context, newExpense CreateExpense) (*Expense, error)
	ListExpenses(ctx context.Context, userID uuid.UUID, transactionDate time.Time) ([]Expense, error)
	GetExpenseByID(ctx context.Context, id uuid.UUID) (*Expense, error)
	DeleteExpense(ctx context.Context, id uuid.UUID) error
}
