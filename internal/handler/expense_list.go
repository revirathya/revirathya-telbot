package handler

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/google/uuid"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/revirathya/revirathya-telbot/internal/domain"
	"github.com/revirathya/revirathya-telbot/internal/handler/subflow"
	"github.com/revirathya/revirathya-telbot/internal/platform/telegram"
	"github.com/revirathya/revirathya-telbot/internal/utils"
)

type ConversationExpenseList struct {
	Telbot        *telegram.TelegramBotHandler
	User          *domain.UserContext
	ChosenExpense *domain.Expense
}

func NewConversationExpenseList(tbh *telegram.TelegramBotHandler, user *domain.UserContext) *ConversationExpenseList {
	return &ConversationExpenseList{
		Telbot: tbh,
		User:   user,
	}
}

// Private
func (c *ConversationExpenseList) newExpenseList(expenseList []domain.Expense) [][]tgbotapi.InlineKeyboardButton {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0)

	// Expenses
	for _, exp := range expenseList {
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(exp.FormatInline(), fmt.Sprintf("expense=%s", exp.ID)),
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

func (c *ConversationExpenseList) deleteExpense(expenseList []domain.Expense, id uuid.UUID) []domain.Expense {
	// Find index
	deleteIdx := slices.IndexFunc(expenseList, func(e domain.Expense) bool {
		return e.ID == id
	})

	// Delete by index
	deletedExpenseList := slices.Delete(expenseList, deleteIdx, deleteIdx+1)
	return deletedExpenseList
}

func (c *ConversationExpenseList) getExpense(expenseList []domain.Expense, expenseID string) domain.Expense {
	pExpenseID, _ := uuid.Parse(expenseID)

	var expense domain.Expense
	for _, e := range expenseList {
		if e.ID == pExpenseID {
			expense = e
			break
		}
	}

	return expense
}

// Public
func (c *ConversationExpenseList) Entry(chatID int64, messageID int) string {
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
	if _, err := c.Telbot.EditMessage(chatID, c.User.InteractionMessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message")
	}

	// Update User Context
	c.User.StartConversation("expense-list", chatID, messageID)
	return "callback:expense-list"
}

func (c *ConversationExpenseList) SelectExpense(expenseList []domain.Expense, data string) string {
	// Set Text
	date := strings.Split(data, "=")[1]

	var text string
	expensesButton := make([][]tgbotapi.InlineKeyboardButton, 0, len(expenseList))

	// Get Expenses
	expenses := make([]domain.Expense, 0)
	for _, e := range expenseList {
		if e.Date == date {
			expenses = append(expenses, e)
		}
	}

	if len(expenses) == 0 {
		// Set Text
		text = fmt.Sprintf("No recorded expenses found for selected date: %s", date)
	} else {
		// Set Text
		text = strings.Join([]string{
			"<b>Expense List</b>",
			fmt.Sprintf("Recorded expenses for %s", date),
		}, "\n\n")

		// Get Expense List buttons
		expensesButton = c.newExpenseList(expenseList)
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
	if _, err := c.Telbot.EditMessage(c.User.InteractionChatID, c.User.InteractionMessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message")
	}

	// Set User Context
	c.User.ChangeState("select-expense")
	return fmt.Sprintf("callback:%s", data)
}

func (c *ConversationExpenseList) PreviewExpense(expenseList []domain.Expense, data string, useChosen bool) string {
	// Get Expense
	var expense domain.Expense
	if useChosen {
		expense = *c.ChosenExpense
	} else {
		expenseID := strings.Split(data, "=")[1]
		expense = c.getExpense(expenseList, expenseID)
		c.ChosenExpense = &expense
	}

	// Set Text
	text := strings.Join([]string{
		"Preview of selected Expense",
		expense.FormatMessage(),
	}, "\n\n")

	// Set Reply Markup
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Delete", "delete"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Back to List", fmt.Sprintf("date=%s", c.ChosenExpense.Date)),
			tgbotapi.NewInlineKeyboardButtonData("Done", "done"),
		),
	)

	// Edit Message
	if _, err := c.Telbot.EditMessage(c.User.InteractionChatID, c.User.InteractionMessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message")
	}

	// Set User Context
	c.User.ChangeState("preview-expense")
	return fmt.Sprintf("callback:%s", data)
}

func (c *ConversationExpenseList) ConfirmationDeleteExpense() string {
	// Set Text
	text := strings.Join([]string{
		"Preview of selected Expense",
		c.ChosenExpense.FormatMessage(),
		"Are you sure want to delete this expense?",
	}, "\n\n")

	// Set Reply Markup
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Yes", "confirmation=yes"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("No", "confirmation=no"),
		),
	)

	// Edit Message
	if _, err := c.Telbot.EditMessage(c.User.InteractionChatID, c.User.InteractionMessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message")
	}

	// Set User Context
	c.User.ChangeState("delete-confirmation")
	return "callback:delete"
}

func (c *ConversationExpenseList) AfterDeleteExpense(expenseList []domain.Expense) string {
	// Remove Expense
	expenseList = c.deleteExpense(expenseList, c.ChosenExpense.ID)

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
	if _, err := c.Telbot.EditMessage(c.User.InteractionChatID, c.User.InteractionMessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message")
	}

	// Set User context
	c.User.ChangeState("after-delete")
	return "callback:after-delete"
}

func (c *ConversationExpenseList) Closure() string {
	// Set Text
	text := "Closing expenses list"

	// Edit Message
	if _, err := c.Telbot.EditMessage(c.User.InteractionChatID, c.User.InteractionMessageID, text, nil); err != nil {
		slog.Error("failed to edit telegram message")
	}

	// Update User Context
	c.ChosenExpense = nil
	c.User.CloseConversation()
	return "callback:done"
}
