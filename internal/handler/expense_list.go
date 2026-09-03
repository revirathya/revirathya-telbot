package handler

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/revirathya/revirathya-telbot/internal/domain"
	"github.com/revirathya/revirathya-telbot/internal/handler/container"
	"github.com/revirathya/revirathya-telbot/internal/handler/subflow"
	"github.com/revirathya/revirathya-telbot/internal/platform/telegram"
	ucExp "github.com/revirathya/revirathya-telbot/internal/usecase/expense"
	ucUtel "github.com/revirathya/revirathya-telbot/internal/usecase/usercontext"

	"github.com/revirathya/revirathya-telbot/internal/utils"
)

type ConversationExpenseList struct {
	telbot *telegram.TelegramBotHandler
	utelUc *ucUtel.UserTelegramContextUsecase
	expUc  *ucExp.ExpenseUsecase
}

func NewConversationExpenseList(
	tbh *telegram.TelegramBotHandler,
	utelUc *ucUtel.UserTelegramContextUsecase,
	expUc *ucExp.ExpenseUsecase,
) *ConversationExpenseList {
	return &ConversationExpenseList{
		telbot: tbh,
		utelUc: utelUc,
		expUc:  expUc,
	}
}

// Private
func (c *ConversationExpenseList) newExpenseList(expenseList []domain.Expense) [][]tgbotapi.InlineKeyboardButton {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0)

	// Expenses
	for _, exp := range expenseList {
		ec := container.NewFromDomain(&exp)
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					ec.FormatInline(),
					fmt.Sprintf("expense=%s", exp.ID),
				),
			))
	}

	// Navigation
	rows = append(
		rows,
		tgbotapi.NewInlineKeyboardRow(
			utils.NewEmptyInlineKeyboardButton("<"),
			utils.NewEmptyInlineKeyboardButton("Page 1/1"),
			utils.NewEmptyInlineKeyboardButton(">"),
		),
	)

	return rows
}

// Public
func (c *ConversationExpenseList) Entry(ctx context.Context, user *domain.UserTelegramContext, chatID int64, messageID int, data string) string {
	// Set Text
	text := strings.Join([]string{
		"<b>List Expense</b>",
		"Please select date:",
	}, "\n\n")

	// Set Reply Markup
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			utils.NewInlineTodayButton("Today", "date="),
			subflow.NewInlineCalendarButton("Pick a Date", nil),
		),
	)

	// Edit Message
	if _, err := c.telbot.EditMessage(chatID, messageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}

	// Update User Context
	if _, err := c.utelUc.StartConversation(ctx, user.UserID, "expense-list", chatID, messageID); err != nil {
		slog.Error("failed to start conversation", "error", err)
	}
	return fmt.Sprintf("callback:%s", data)
}

func (c *ConversationExpenseList) SelectExpense(ctx context.Context, user *domain.UserTelegramContext, data string) string {
	trxDate, err := time.Parse("2006-01-02", strings.Split(data, "=")[1])
	if err != nil {
		slog.Error("invalid transaction_dt Date format", "error", err)
	}

	// Get Expenses
	expenses, err := c.expUc.ListExpenses(ctx, user.UserID, trxDate)
	if err != nil {
		slog.Error("failed to list expenses", "error", err)
	}

	expensesButton := make([][]tgbotapi.InlineKeyboardButton, 0, len(expenses))

	var text string
	if len(expenses) == 0 {
		// Set Text
		text = fmt.Sprintf("No recorded expenses found for selected date: %s", trxDate.Format("2006-01-02"))
	} else {
		// Set Text
		text = strings.Join([]string{
			"<b>Expense List</b>",
			fmt.Sprintf("Recorded expenses for %s", trxDate.Format("2006-01-02")),
		}, "\n\n")

		// Get Expense List buttons
		expensesButton = c.newExpenseList(expenses)
	}

	utilityButton := tgbotapi.NewInlineKeyboardRow(
		subflow.NewInlineCalendarButton("Pick Other Date", nil),
		tgbotapi.NewInlineKeyboardButtonData("Done", "done"),
	)

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	rows = append(rows, expensesButton...)
	rows = append(rows, utilityButton)

	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	// Edit Message
	if _, err := c.telbot.EditMessage(*user.InteractionChatID, *user.InteractionMessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}

	// Set User Context
	if _, err := c.utelUc.ChangeState(ctx, user.UserID, "select-expense"); err != nil {
		slog.Error("failed to edit change state", "error", err)
	}
	return fmt.Sprintf("callback:%s", data)
}

func (c *ConversationExpenseList) PreviewExpense(ctx context.Context, user *domain.UserTelegramContext, data string) string {
	// Get Expense
	expenseID, _ := uuid.Parse(strings.Split(data, "=")[1])
	exp, _ := c.expUc.GetExpenseByID(ctx, expenseID)
	ec := container.NewFromDomain(exp)

	// Set Text
	text := strings.Join([]string{
		"Preview of selected Expense",
		ec.FormatDetail(),
	}, "\n\n")

	// Set Reply Markup
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Delete",
				fmt.Sprintf("delete=%s", expenseID),
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Back to List",
				fmt.Sprintf("date=%s", exp.TransactionDate.Format("2006-01-02")),
			),
			tgbotapi.NewInlineKeyboardButtonData("Done", "done"),
		),
	)

	// Edit Message
	if _, err := c.telbot.EditMessage(*user.InteractionChatID, *user.InteractionMessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}

	// Set User Context
	if _, err := c.utelUc.ChangeState(ctx, user.UserID, "preview-expense"); err != nil {
		slog.Error("failed to change state", "error", err)
	}
	return fmt.Sprintf("callback:%s", data)
}

func (c *ConversationExpenseList) ConfirmationDeleteExpense(ctx context.Context, user *domain.UserTelegramContext, data string) string {
	// Get Expense
	expenseID, err := uuid.Parse(strings.Split(data, "=")[1])
	if err != nil {
		slog.Error("invalid expense_id UUID format", "error", err)
	}

	exp, err := c.expUc.GetExpenseByID(ctx, expenseID)
	if err != nil {
		slog.Error("failed to fetch expense", "error", err)
	}

	// Create Container
	ec := container.NewFromDomain(exp)

	// Set Text
	text := strings.Join([]string{
		"Preview of selected Expense",
		ec.FormatDetail(),
		"Are you sure want to delete this expense?",
	}, "\n\n")

	// Set Reply Markup
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Yes", fmt.Sprintf("confirmation=%s", expenseID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("No", "confirmation=no"),
		),
	)

	// Edit Message
	if _, err := c.telbot.EditMessage(*user.InteractionChatID, *user.InteractionMessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}

	// Set User Context
	if _, err := c.utelUc.ChangeState(ctx, user.UserID, "delete-confirmation"); err != nil {
		slog.Error("failed to change state", "error", err)
	}
	return fmt.Sprintf("callback:%s", data)
}

func (c *ConversationExpenseList) AfterDeleteExpense(ctx context.Context, user *domain.UserTelegramContext, data string) string {
	// Delete Expense
	expenseID, err := uuid.Parse(strings.Split(data, "=")[1])
	if err != nil {
		slog.Error("invalid expense_id UUID format", "error", err)
	}

	if err := c.expUc.DeleteExpense(ctx, expenseID); err != nil {
		slog.Error("failed to delete expense", "error", err)
	}

	// Set Text
	text := "Expense successfully deleted"

	// Set Inline Reply Markup
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			subflow.NewInlineCalendarButton("Pick Other Date", nil),
			tgbotapi.NewInlineKeyboardButtonData("Done", "done"),
		),
	)

	// Edit Message
	if _, err := c.telbot.EditMessage(*user.InteractionChatID, *user.InteractionMessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}

	// Set User context
	if _, err := c.utelUc.ChangeState(ctx, user.UserID, "after-delete"); err != nil {
		slog.Error("failed to change state", "error", err)
	}
	return fmt.Sprintf("callback:%s", data)
}

func (c *ConversationExpenseList) Closure(ctx context.Context, user *domain.UserTelegramContext, data string) string {
	// Set Text
	text := "Closing expenses list"

	// Edit Message
	if _, err := c.telbot.EditMessage(*user.InteractionChatID, *user.InteractionMessageID, text, nil); err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}

	// Update User Context
	if _, err := c.utelUc.CloseConversation(ctx, user.UserID); err != nil {
		slog.Error("failed to close conversation", "error", err)
	}
	return fmt.Sprintf("callback:%s", data)
}
