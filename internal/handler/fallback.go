package handler

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/revirathya/revirathya-telbot/internal/domain"
	"github.com/revirathya/revirathya-telbot/internal/platform/telegram"
	ucUtel "github.com/revirathya/revirathya-telbot/internal/usecase/usercontext"
)

func CommandCancel(tbh *telegram.TelegramBotHandler, user *domain.UserTelegramContext, ctx context.Context, utelUc *ucUtel.UserTelegramContextUsecase, chatID int64, command string) string {
	// Set Text
	var text string
	if user.Conversation != nil {
		text = "Cancel conversation"
	} else {
		text = "There is no active conversation"
	}

	// Send Message
	if _, err := tbh.SendMessage(chatID, text, nil); err != nil {
		slog.Error("failed to send telegram message", "error", err)
	}

	// Clear Context
	if _, err := utelUc.CloseConversation(ctx, user.UserID); err != nil {
		slog.Error("failed to send telegram message", "error", err)
	}
	return fmt.Sprintf("command:%s", command)
}
