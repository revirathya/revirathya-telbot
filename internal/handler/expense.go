package handler

import (
	"context"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/revirathya/revirathya-telbot/internal/domain"
	"github.com/revirathya/revirathya-telbot/internal/platform/telegram"
	ucUtel "github.com/revirathya/revirathya-telbot/internal/usecase/usercontext"
)

type ConversationExpense struct {
	telbot *telegram.TelegramBotHandler
	utelUc *ucUtel.UserTelegramContextUsecase
}

func NewConversationExpense(tbh *telegram.TelegramBotHandler, utelUc *ucUtel.UserTelegramContextUsecase) *ConversationExpense {
	return &ConversationExpense{
		telbot: tbh,
		utelUc: utelUc,
	}
}

type Option struct {
	Label string
	Value string
}

func (c *ConversationExpense) Entry(ctx context.Context, user *domain.UserTelegramContext, chatID int64) string {
	// Set Text
	text := "Please choose <b>expenses</b> actions:"

	// Set Reply Markup
	options := []Option{
		{Label: "Record", Value: "expense-record"},
		{Label: "List", Value: "expense-list"},
		{Label: "Summary", Value: "expense-summary"},
	}
	inlineKeyboard := make([][]tgbotapi.InlineKeyboardButton, 0, len(options))
	for _, opt := range options {
		inlineKeyboard = append(inlineKeyboard, []tgbotapi.InlineKeyboardButton{{Text: opt.Label, CallbackData: &opt.Value}})
	}
	replyMarkup := &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: inlineKeyboard}

	// Send Message
	upt, err := c.telbot.SendMessage(chatID, text, replyMarkup)
	if err != nil {
		slog.Error("failed to edit telegram message")
	}

	// Update user context
	c.utelUc.StartConversation(ctx, user.UserID, "expense", chatID, upt.MessageID)
	return "command:expense"
}
