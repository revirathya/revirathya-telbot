package usecase

import (
	// Internal
	"github.com/revirathya/revirathya-telbot/internal/domain"

	// Usecase
	"github.com/revirathya/revirathya-telbot/internal/usecase/category"
	"github.com/revirathya/revirathya-telbot/internal/usecase/expense"
	"github.com/revirathya/revirathya-telbot/internal/usecase/usertelegramcontext"
)

// Bootstraping usecase
type Usecase struct {
	Category            *category.CategoryUsecase
	Expense             *expense.ExpenseUsecase
	UserTelegramContext *usertelegramcontext.UserTelegramContextUsecase
}

func BootstrapUsecase(repo domain.Repository) Usecase {
	return Usecase{
		Category:            category.NewCategoryUsecase(repo.Category()),
		Expense:             expense.NewExpenseUsecase(repo.Expense()),
		UserTelegramContext: usertelegramcontext.NewUserTelegramContextUsecase(repo.UserTelegramContext()),
	}
}
