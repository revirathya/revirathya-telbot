package domain

import (
	"context"

	"github.com/google/uuid"
)

// Entity: User Telegram Context
type UserTelegramContext struct {
	UserID               uuid.UUID
	TelegramChatID       *string
	Conversation         *string
	State                *string
	InteractionChatID    *int64
	InteractionMessageID *int
}

// Repository: User Telegram Context
type UserTelegramContextRepository interface {
	GetUserContextByTelegramID(ctx context.Context, telegramChatID string) (*UserTelegramContext, error)
	UpdateUserContext(ctx context.Context, userID uuid.UUID, conversation *string, state *string, interactionChatID *int64, interactionMessageID *int) (*UserTelegramContext, error)
}
