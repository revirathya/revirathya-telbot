package repository

import (
	"github.com/revirathya/revirathya-telbot/internal/domain"
	"github.com/revirathya/revirathya-telbot/internal/platform/grpc/clients"
	"github.com/revirathya/revirathya-telbot/internal/repository/grpc/category"
	"github.com/revirathya/revirathya-telbot/internal/repository/grpc/expense"
	"github.com/revirathya/revirathya-telbot/internal/repository/grpc/usertelegramcontext"
)

// Bootstraping repository
type repositoryImpl struct {
	repoCategory            domain.CategoryRepository
	repoExpense             domain.ExpenseRepository
	repoUserTelegramContext domain.UserTelegramContextRepository
}

func BootstrapRepository(finance *clients.FinanceClient) domain.Repository {
	return &repositoryImpl{
		repoCategory:            category.NewCategoryRepositoryImpl(finance),
		repoExpense:             expense.NewExpenseRepositoryImpl(finance),
		repoUserTelegramContext: usertelegramcontext.NewUserTelegramContextRepositoryImpl(finance),
	}
}

func (r *repositoryImpl) Category() domain.CategoryRepository {
	return r.repoCategory
}
func (r *repositoryImpl) Expense() domain.ExpenseRepository {
	return r.repoExpense
}
func (r *repositoryImpl) UserTelegramContext() domain.UserTelegramContextRepository {
	return r.repoUserTelegramContext
}
