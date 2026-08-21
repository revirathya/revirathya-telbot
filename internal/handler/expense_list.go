package handler

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/revirathya/revirathya-telbot/infra/telegram"
	"github.com/revirathya/revirathya-telbot/internal/domain"
)

type ConversationExpenseList struct {
	Telbot *telegram.TelegramBotHandler
	User   *domain.UserContext
}

func NewConversationExpenseList(tbh *telegram.TelegramBotHandler, user *domain.UserContext) *ConversationExpenseList {
	return &ConversationExpenseList{
		Telbot: tbh,
		User:   user,
	}
}

func (c *ConversationExpenseList) Entry(chatID int64, messageID int) string {
	// Set Text
	text := strings.Join([]string{
		"<b>List Expense</b>",
		"Please select date:",
	}, "\n\n")

	// Set Reply Markup
	todayDt := time.Now()
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Today", fmt.Sprintf("date=%s", todayDt.Format("2006-01-02"))),
			tgbotapi.NewInlineKeyboardButtonData("Pick a Date", fmt.Sprintf("calendar=%s", todayDt.Format("2006-01"))),
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

func (c *ConversationExpenseList) Closure(expenseList []domain.Expense, data string) string {
	// Set Text
	date := strings.Split(data, "=")[1]

	var text string
	var ptrInlineKeyboardMarkup *tgbotapi.InlineKeyboardMarkup

	if len(expenseList) == 0 {
		// Set Text
		text = fmt.Sprintf("No recorded expenses found for selected date: %s", date)

		// Set Reply Markup
		ptrInlineKeyboardMarkup = nil
	} else {
		// Set Text
		text = strings.Join([]string{
			"<b>Expense List</b>",
			fmt.Sprintf("Recorded expenses for %s", date),
		}, "\n\n")

		// Set Reply Markup
		rows := make([][]tgbotapi.InlineKeyboardButton, 0)
		for _, exp := range expenseList {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(exp.FormatInline(), fmt.Sprintf("expense=%d", exp.CategoryID)),
			))
		}
		inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(rows...)
		ptrInlineKeyboardMarkup = &inlineKeyboardMarkup
	}

	// Edit Message
	if _, err := c.Telbot.EditMessage(c.User.InteractionChatID, c.User.InteractionMessageID, text, ptrInlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message")
	}

	// Update User Context
	c.User.CloseConversation()
	return fmt.Sprintf("callback:%s", data)
}
