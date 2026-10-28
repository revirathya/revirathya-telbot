package bot

import (
	"github.com/revirathya/revirathya-telbot/internal/bot/handler"
	"github.com/revirathya/revirathya-telbot/internal/platform/telegram"
	"github.com/revirathya/revirathya-telbot/internal/usecase"
)

type BotHandler struct {
	Expense        *handler.ConversationExpense
	ExpenseRecord  *handler.ConversationExpenseRecord
	ExpenseList    *handler.ConversationExpenseList
	ExpenseSummary *handler.ConversationExpenseSummary
	CommandCancel  *handler.CommandCancel
}

func BootstrapBotHandler(tbh *telegram.TelegramBotHandler, uc usecase.Usecase) *BotHandler {
	return &BotHandler{
		Expense:        handler.NewConversationExpense(tbh, uc.UserTelegramContext),
		ExpenseRecord:  handler.NewConversationExpenseRecord(tbh, uc.UserTelegramContext, uc.Category, uc.Expense),
		ExpenseList:    handler.NewConversationExpenseList(tbh, uc.UserTelegramContext, uc.Expense),
		ExpenseSummary: handler.NewConversationExpenseSummary(tbh, uc.UserTelegramContext, uc.Expense),
		CommandCancel:  handler.NewCommandCancel(tbh, uc.UserTelegramContext),
	}
}
