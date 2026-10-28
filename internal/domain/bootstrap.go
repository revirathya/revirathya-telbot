package domain

type Repository interface {
	Category() CategoryRepository
	Expense() ExpenseRepository
	UserTelegramContext() UserTelegramContextRepository
}
