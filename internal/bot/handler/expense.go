package handler

import (
	"context"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	// Internal
	"github.com/revirathya/revirathya-telbot/internal/bot/updates"
	"github.com/revirathya/revirathya-telbot/internal/platform/telegram"
)

type ConversationExpense struct {
	tbh                 *telegram.TelegramBotHandler
	userTelegramContext userTelegramContextUsecase
}

func NewConversationExpense(tbh *telegram.TelegramBotHandler, userTelegramContext userTelegramContextUsecase) *ConversationExpense {
	return &ConversationExpense{
		tbh:                 tbh,
		userTelegramContext: userTelegramContext,
	}
}

type Option struct {
	Label string
	Value string
}

func (c *ConversationExpense) Entry(ctx context.Context, upt updates.UserUpdate) string {
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
	msg, err := c.tbh.SendMessage(upt.Incoming.ChatID, text, replyMarkup)
	if err != nil {
		slog.Error("failed to edit telegram message")
	}

	// Update user context
	if _, err := c.userTelegramContext.StartConversation(ctx, "expense", msg.Chat.ID, msg.MessageID); err != nil {
		slog.Error("failed to start conversation")
	}
	return "command:expense"
}
