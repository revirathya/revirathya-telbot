package updates

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Get User
type UserLookup func(context.Context, tgbotapi.Update) (UserUpdate, error)

func GetUser(ctx context.Context, tgUpt tgbotapi.Update, userLookup UserLookup) (UserUpdate, error) {
	// Telegram ID
	user, err := userLookup(ctx, tgUpt)
	if err != nil {
		return UserUpdate{}, err
	}

	return user, nil
}
