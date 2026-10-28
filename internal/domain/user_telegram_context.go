package domain

import (
	"context"
)

// Entity
type UserTelegramContext struct {
	TelegramChatID       *string
	Conversation         *string
	State                *string
	InteractionChatID    *int64
	InteractionMessageID *int
}

// Patch
type UserTelegramContextPatch struct {
	Conversation         *string
	State                *string
	InteractionChatID    *int64
	InteractionMessageID *int
}

// Repository: User Telegram Context
type UserTelegramContextRepository interface {
	GetUserContext(ctx context.Context) (*UserTelegramContext, error)
	UpdateUserContext(ctx context.Context, patch UserTelegramContextPatch) (*UserTelegramContext, error)
}
