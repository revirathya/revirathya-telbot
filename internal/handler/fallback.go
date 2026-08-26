package handler

import (
	"log/slog"

	"github.com/revirathya/revirathya-telbot/internal/domain"
	"github.com/revirathya/revirathya-telbot/internal/platform/telegram"
)

func CommandCancel(tbh *telegram.TelegramBotHandler, user *domain.UserContext, chatID int64) string {
	// Set Text
	var text string
	if user.Conversation != "root" {
		text = "Cancel conversation"
	} else {
		text = "There is no active conversation"
	}

	// Send Message
	if _, err := tbh.SendMessage(chatID, text, nil); err != nil {
		slog.Error("failed to send telegram message")
	}

	// Clear Context
	user.CloseConversation()
	return "command:cancel"
}
