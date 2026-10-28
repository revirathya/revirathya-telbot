package handler

import (
	"context"
	"log/slog"

	// Internal
	"github.com/revirathya/revirathya-telbot/internal/bot/updates"
	"github.com/revirathya/revirathya-telbot/internal/platform/telegram"
)

type CommandCancel struct {
	tbh                 *telegram.TelegramBotHandler
	userTelegramContext userTelegramContextUsecase
}

func NewCommandCancel(tbh *telegram.TelegramBotHandler, userTelegramContext userTelegramContextUsecase) *CommandCancel {
	return &CommandCancel{
		tbh:                 tbh,
		userTelegramContext: userTelegramContext,
	}
}

func (c *CommandCancel) Command(ctx context.Context, upt updates.UserUpdate) {
	// Set Text
	var text string
	if upt.Interaction.Conversation == "root" {
		// Send no active conversation to be cancelled
		text = "There is no active conversation"
		if _, err := c.tbh.SendMessage(upt.Interaction.ChatID, text, nil); err != nil {
			slog.Error("failed to send telegram message", "error", err)
		}
		return
	}

	// Send cancel conversation
	text = "Cancel conversation"
	if _, err := c.tbh.SendMessage(upt.Interaction.ChatID, text, nil); err != nil {
		slog.Error("failed to send telegram message", "error", err)
	}

	// Clear context
	if _, err := c.userTelegramContext.CloseConversation(ctx); err != nil {
		slog.Error("failed to send telegram message", "error", err)
	}
}
