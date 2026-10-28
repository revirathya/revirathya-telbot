package handler

import (
	"context"
	"time"

	"github.com/google/uuid"

	// Internal
	"github.com/revirathya/revirathya-telbot/internal/domain"
)

// Category
type categoryUsecase interface {
	ListCategories(ctx context.Context) ([]domain.Category, error)
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*domain.Category, error)
}

// Expense
type expenseUsecase interface {
	CreateExpense(ctx context.Context, form domain.ExpenseForm, requestID uuid.UUID) (*domain.Expense, error)
	ListExpenses(ctx context.Context, transactionDate time.Time) ([]domain.Expense, error)
	GetExpenseByID(ctx context.Context, id uuid.UUID) (*domain.Expense, error)
	DeleteExpense(ctx context.Context, id uuid.UUID) error
}

// User Telegram Context
type userTelegramContextUsecase interface {
	StartConversation(ctx context.Context, conversation string, chatID int64, messageID int) (*domain.UserTelegramContext, error)
	ChangeState(ctx context.Context, state string) (*domain.UserTelegramContext, error)
	CloseConversation(ctx context.Context) (*domain.UserTelegramContext, error)
}
