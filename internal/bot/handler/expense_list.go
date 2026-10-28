package handler

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	// Internal
	"github.com/revirathya/revirathya-telbot/internal/bot/container"
	"github.com/revirathya/revirathya-telbot/internal/bot/subflow"
	"github.com/revirathya/revirathya-telbot/internal/bot/updates"
	"github.com/revirathya/revirathya-telbot/internal/domain"
	"github.com/revirathya/revirathya-telbot/internal/platform/telegram"
	"github.com/revirathya/revirathya-telbot/internal/utils"
)

type ConversationExpenseList struct {
	tbh                 *telegram.TelegramBotHandler
	userTelegramContext userTelegramContextUsecase
	expense             expenseUsecase
}

func NewConversationExpenseList(
	tbh *telegram.TelegramBotHandler,
	userTelegramContext userTelegramContextUsecase,
	expense expenseUsecase,
) *ConversationExpenseList {
	return &ConversationExpenseList{
		tbh:                 tbh,
		userTelegramContext: userTelegramContext,
		expense:             expense,
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
func (c *ConversationExpenseList) Entry(ctx context.Context, upt updates.UserUpdate) string {
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
	msg, err := c.tbh.EditMessage(upt.Interaction.ChatID, upt.Interaction.MessageID, text, &inlineKeyboardMarkup)
	if err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}

	// Update User Context
	if _, err := c.userTelegramContext.StartConversation(ctx, "expense-list", msg.Chat.ID, msg.MessageID); err != nil {
		slog.Error("failed to start conversation", "error", err)
	}
	return fmt.Sprintf("callback:%s", upt.Incoming.Data)
}

func (c *ConversationExpenseList) SelectExpense(ctx context.Context, upt updates.UserUpdate) string {
	trxDate, err := time.Parse("2006-01-02", upt.Incoming.Data)
	if err != nil {
		slog.Error("invalid transaction_dt Date format", "error", err)
	}

	// Get Expenses
	expenses, err := c.expense.ListExpenses(ctx, trxDate)
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
	if _, err := c.tbh.EditMessage(upt.Interaction.ChatID, upt.Interaction.MessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}

	// Set User Context
	if _, err := c.userTelegramContext.ChangeState(ctx, "select-expense"); err != nil {
		slog.Error("failed to edit change state", "error", err)
	}
	return fmt.Sprintf("callback:%s", upt.Incoming.Data)
}

func (c *ConversationExpenseList) PreviewExpense(ctx context.Context, upt updates.UserUpdate) string {
	// Get Expense
	expID, err := uuid.Parse(upt.Incoming.Data)
	if err != nil {
		slog.Error("invalid expense_id UUID format", "error", err)
	}

	exp, err := c.expense.GetExpenseByID(ctx, expID)
	if err != nil {
		slog.Error("failed to fetch expense", "error", err)
	}

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
				fmt.Sprintf("delete=%s", expID),
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Back to List",
				fmt.Sprintf("date=%s", exp.TransactionDt.Format("2006-01-02")),
			),
			tgbotapi.NewInlineKeyboardButtonData("Done", "done"),
		),
	)

	// Edit Message
	if _, err := c.tbh.EditMessage(upt.Interaction.ChatID, upt.Interaction.MessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}

	// Set User Context
	if _, err := c.userTelegramContext.ChangeState(ctx, "preview-expense"); err != nil {
		slog.Error("failed to change state", "error", err)
	}
	return fmt.Sprintf("callback:%s", upt.Incoming.Data)
}

func (c *ConversationExpenseList) ConfirmationDeleteExpense(ctx context.Context, upt updates.UserUpdate) string {
	// Get Expense
	expID, err := uuid.Parse(upt.Incoming.Data)
	if err != nil {
		slog.Error("invalid expense_id UUID format", "error", err)
	}

	exp, err := c.expense.GetExpenseByID(ctx, expID)
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
			tgbotapi.NewInlineKeyboardButtonData("Yes", fmt.Sprintf("confirmation=%s", expID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("No", "confirmation=no"),
		),
	)

	// Edit Message
	if _, err := c.tbh.EditMessage(upt.Interaction.ChatID, upt.Interaction.MessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}

	// Set User Context
	if _, err := c.userTelegramContext.ChangeState(ctx, "delete-confirmation"); err != nil {
		slog.Error("failed to change state", "error", err)
	}
	return fmt.Sprintf("callback:%s", upt.Incoming.Data)
}

func (c *ConversationExpenseList) AfterDeleteExpense(ctx context.Context, upt updates.UserUpdate) string {
	// Delete Expense
	expID, err := uuid.Parse(upt.Incoming.Data)
	if err != nil {
		slog.Error("invalid expense_id UUID format", "error", err)
	}

	if err := c.expense.DeleteExpense(ctx, expID); err != nil {
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
	if _, err := c.tbh.EditMessage(upt.Interaction.ChatID, upt.Interaction.MessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}

	// Set User context
	if _, err := c.userTelegramContext.ChangeState(ctx, "after-delete"); err != nil {
		slog.Error("failed to change state", "error", err)
	}
	return fmt.Sprintf("callback:%s", upt.Incoming.Data)
}

func (c *ConversationExpenseList) Closure(ctx context.Context, upt updates.UserUpdate) string {
	// Set Text
	text := "Closing expenses list"

	// Edit Message
	if _, err := c.tbh.EditMessage(upt.Interaction.ChatID, upt.Interaction.MessageID, text, nil); err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}

	// Update User Context
	if _, err := c.userTelegramContext.CloseConversation(ctx); err != nil {
		slog.Error("failed to close conversation", "error", err)
	}
	return fmt.Sprintf("callback:%s", upt.Incoming.Data)
}
