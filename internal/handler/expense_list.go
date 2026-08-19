package handler

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/revirathya/revirathya-telbot/infra/telegram"
	"github.com/revirathya/revirathya-telbot/internal/domain"
	"github.com/revirathya/revirathya-telbot/internal/utils"
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
	inlineKeyboardMarkup := &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{{
				Text:         "Today",
				CallbackData: utils.ToStrPtr(fmt.Sprintf("date=%s", time.Now().Format("2006-01-02"))),
			}},
			{{
				Text:         "Pick a Date",
				CallbackData: utils.ToStrPtr("calendar"),
			}},
		},
	}

	// Edit Message
	if _, err := c.Telbot.EditMessage(chatID, c.User.InteractionMessageID, text, inlineKeyboardMarkup); err != nil {
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
	var inlineKeyboardMarkup *tgbotapi.InlineKeyboardMarkup

	if len(expenseList) == 0 {
		// Set Text
		text = fmt.Sprintf("No recorded expenses found for selected date: %s", date)

		// Set Reply Markup
		inlineKeyboardMarkup = nil
	} else {
		// Set Text
		text = strings.Join([]string{
			"<b>Expense List</b>",
			fmt.Sprintf("Recorded expenses for %s", date),
		}, "\n\n")

		// Set Reply Markup
		inlineKeyboard := make([][]tgbotapi.InlineKeyboardButton, 0)
		for _, exp := range expenseList {
			value := fmt.Sprintf("expense=%d", exp.CategoryID) // TMP
			inlineKeyboard = append(inlineKeyboard, []tgbotapi.InlineKeyboardButton{{Text: exp.FormatInline(), CallbackData: &value}})
		}
		inlineKeyboardMarkup = &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: inlineKeyboard}
	}

	// Edit Message
	if _, err := c.Telbot.EditMessage(c.User.InteractionChatID, c.User.InteractionMessageID, text, inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message")
	}

	// Update User Context
	c.User.CloseConversation()
	return fmt.Sprintf("callback:%s", data)
}
