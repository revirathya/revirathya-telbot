package handler

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/revirathya/revirathya-telbot/infra/telegram"
	"github.com/revirathya/revirathya-telbot/internal/domain"
	"github.com/revirathya/revirathya-telbot/internal/handler/subflow"
	"github.com/revirathya/revirathya-telbot/internal/utils"
)

type ConversationExpenseRecord struct {
	Telbot        *telegram.TelegramBotHandler
	User          *domain.UserContext
	CreateExpense *domain.Expense
}

func NewConversationExpenseRecord(tbh *telegram.TelegramBotHandler, user *domain.UserContext) *ConversationExpenseRecord {
	return &ConversationExpenseRecord{
		Telbot:        tbh,
		User:          user,
		CreateExpense: &domain.Expense{},
	}
}

func (c *ConversationExpenseRecord) Entry(chatID int64, messageID int) string {
	// Set Text
	text := strings.Join([]string{
		"<b>Record Expense</b>",
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
	c.User.StartConversation("expense-record", chatID, messageID)
	return "callback:expense-record"
}

func (c *ConversationExpenseRecord) InteractionSelectCategory(categories []domain.Category, data string) string {
	// Update Create Expense Context
	inpDate := strings.Split(data, "=")[1]
	c.CreateExpense.Date = inpDate

	// Set Text
	text := strings.Join([]string{
		"<b>Record Expense</b>",
		c.CreateExpense.FormatMessage(),
		"Please select category:",
	}, "\n\n")

	// Set Reply Markup
	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, ctg := range categories {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ctg.Label, fmt.Sprintf("category=%s", ctg.ID)),
		))
	}
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	// Edit Message
	if _, err := c.Telbot.EditMessage(c.User.InteractionChatID, c.User.InteractionMessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message")
	}

	// Update User Context
	c.User.ChangeState("select-category")
	return fmt.Sprintf("callback:%s", data)
}

func (c *ConversationExpenseRecord) InteractionInputDescription(categories []domain.Category, data string) string {
	// Update Create Expense Context
	categoryID := strings.Split(data, "=")[1]
	c.CreateExpense.SetCategory(categoryID, categories)

	// Set text
	text := strings.Join([]string{
		"<b>Record Expense</b>",
		c.CreateExpense.FormatMessage(),
		"Please input description:",
	}, "\n\n")

	// Edit Message
	if _, err := c.Telbot.EditMessage(c.User.InteractionChatID, c.User.InteractionMessageID, text, nil); err != nil {
		slog.Error("failed to edit telegram message")
	}

	// Update User Context
	c.User.ChangeState("input-description")
	return fmt.Sprintf("callback:%s", data)
}

func (c *ConversationExpenseRecord) InteractionInputAmount(text string) string {
	// Update Create Expense Context
	c.CreateExpense.Description = text

	// Set Text
	text = strings.Join([]string{
		"<b>Record Expense</b>",
		c.CreateExpense.FormatMessage(),
		"Please input amount:",
	}, "\n\n")

	// Set Reply Markup
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Input Amount", "numpad=0"),
		),
	)

	// Edit Message
	if _, err := c.Telbot.EditMessage(c.User.InteractionChatID, c.User.InteractionMessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message")
	}

	// Update User Context
	c.User.ChangeState("input-amount")
	return "text:description"
}

func (c *ConversationExpenseRecord) InteractionSelectConfirmation(data string) string {
	// Update Create Expense Context
	amountStr := strings.Split(data, "=")[1]
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		slog.Error("failed to parse int", "amount", amountStr)
	}
	c.CreateExpense.Amount = amount

	// Set Text
	text := strings.Join([]string{
		"<b>Record Expense</b>",
		c.CreateExpense.FormatMessage(),
		"Do you want record this expense?",
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
	c.User.ChangeState("select-confirmation")
	return fmt.Sprintf("callback:%s", data)
}

func (c *ConversationExpenseRecord) Closure(expenseList *[]domain.Expense, data string) string {
	var text string
	confirmation := strings.Split(data, "=")[1]

	switch confirmation {
	case "yes":
		// Set Text
		text = strings.Join([]string{
			"Successfully record following expense:",
			c.CreateExpense.FormatMessage(),
		}, "\n\n")

		// Record expenses to storage
		expense := domain.CreateExpense(
			c.CreateExpense.Date,
			c.CreateExpense.CategoryID,
			c.CreateExpense.CategoryLabel,
			c.CreateExpense.Description,
			c.CreateExpense.Amount,
		)
		*expenseList = append(*expenseList, expense)
	case "no":
		// Set Text
		text = "Cancel record new expense"
	}

	// Edit Message
	if _, err := c.Telbot.EditMessage(c.User.InteractionChatID, c.User.InteractionMessageID, text, nil); err != nil {
		slog.Error("failed to edit telegram message")
	}

	// Clear All Context
	c.CreateExpense.Clear()
	c.User.CloseConversation()
	return fmt.Sprintf("callback:%s", data)
}
