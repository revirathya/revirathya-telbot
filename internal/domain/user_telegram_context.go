package domain

import (
	"context"

	"github.com/revirathya/revirathya-telbot/internal/domain/validation"
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

func (p *UserTelegramContextPatch) Validate() error {
	// Conversation
	if p.Conversation != nil {
		description, err := validation.ValidateString("conversation", *p.Conversation, 100, false)
		if err != nil {
			return err
		}
		p.Conversation = &description
	}

	// State
	if p.State != nil {
		state, err := validation.ValidateString("state", *p.State, 100, false)
		if err != nil {
			return err
		}
		p.State = &state
	}

	// Interaction Chat ID
	if err := validation.ValidateInt64("interaction_chat_id", *p.InteractionChatID, false); err != nil {
		return err
	}

	// Interaction Message ID
	if err := validation.ValidateInt("interaction_message_id", *p.InteractionMessageID, false); err != nil {
		return err
	}
	return nil
}

// Repository: User Telegram Context
type UserTelegramContextRepository interface {
	GetUserContext(ctx context.Context) (*UserTelegramContext, error)
	UpdateUserContext(ctx context.Context, patch UserTelegramContextPatch) (*UserTelegramContext, error)
}
