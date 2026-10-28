package auth

import (
	"context"

	// Internal
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	errDom "github.com/revirathya/revirathya-telbot/internal/domain/error"
)

type telegramIDKey struct{}

func extractID(tgUpt tgbotapi.Update) int64 {
	var telegramID int64 = 0
	if tgUpt.Message != nil {
		telegramID = tgUpt.Message.From.ID
	} else if tgUpt.CallbackQuery != nil {
		telegramID = tgUpt.CallbackQuery.From.ID
	}
	return telegramID
}

// WithTelegramID returns a copy of ctx carrying telegramID.
func WithTelegramID(ctx context.Context, tgUpt tgbotapi.Update) context.Context {
	telegramID := extractID(tgUpt)
	return context.WithValue(ctx, telegramIDKey{}, telegramID)
}

func TelegramIDFrom(ctx context.Context) (int64, error) {
	telegramID, ok := ctx.Value(telegramIDKey{}).(int64)
	if !ok {
		return 0, errDom.ErrUnauthenticated
	}
	return telegramID, nil
}
