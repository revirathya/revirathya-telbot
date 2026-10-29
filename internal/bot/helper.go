package bot

import (
	"context"
	"strings"

	// Internal
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/revirathya/revirathya-telbot/internal/bot/updates"
	"github.com/revirathya/revirathya-telbot/internal/domain"
	"github.com/revirathya/revirathya-telbot/internal/usecase"
)

func coalesceNil[T any](v *T, df T) T {
	if v == nil {
		return df
	}
	return *v
}

func parseCallback(tgUpt tgbotapi.Update, user *domain.UserTelegramContext) updates.ConversationContext {
	// Early pruning
	if tgUpt.CallbackQuery == nil {
		return updates.ConversationContext{}
	}

	var key string
	var data string
	var chatID int64
	var messageID int

	// Parse callback key-value
	parsed := strings.Split(tgUpt.CallbackQuery.Data, "=")
	if len(parsed) > 1 {
		key = strings.Split(tgUpt.CallbackQuery.Data, "=")[0]
		data = strings.Split(tgUpt.CallbackQuery.Data, "=")[1]
	} else {
		key = tgUpt.CallbackData()
		data = tgUpt.CallbackData()
	}

	// Parse chat and message ID
	if user != nil {
		chatID = *user.InteractionChatID
		messageID = *user.InteractionMessageID
	} else {
		chatID = 0
		messageID = 0
	}
	return updates.ConversationContext{
		ChatID:    chatID,
		MessageID: messageID,
		Mode:      "callback",
		Key:       key,
		Data:      data,
	}
}

func parseCommand(tgUpt tgbotapi.Update) updates.ConversationContext {
	// Early pruning
	if (tgUpt.Message == nil) || (tgUpt.Message.Text[0] != '/') {
		return updates.ConversationContext{}
	}

	return updates.ConversationContext{
		ChatID:    tgUpt.Message.Chat.ID,
		MessageID: tgUpt.Message.MessageID,
		Mode:      "command",
		Key:       "command",
		Data:      tgUpt.Message.Text[1:],
	}
}

func parseText(tgUpt tgbotapi.Update) updates.ConversationContext {
	// Early pruning
	if tgUpt.Message == nil {
		return updates.ConversationContext{}
	}

	return updates.ConversationContext{
		ChatID:    tgUpt.Message.Chat.ID,
		MessageID: tgUpt.Message.MessageID,
		Mode:      "text",
		Key:       "text",
		Data:      tgUpt.Message.Text,
	}
}

func UserLookupFunc(usecase usecase.Usecase) func(ctx context.Context, tgUpt tgbotapi.Update) (updates.UserUpdate, error) {
	return func(ctx context.Context, tgUpt tgbotapi.Update) (updates.UserUpdate, error) {
		// Get user contexts
		user, err := usecase.UserTelegramContext.GetUserContext(ctx)
		if err != nil {
			return updates.UserUpdate{}, err
		}

		// Check mode
		var cnvCtx updates.ConversationContext
		switch {
		case tgUpt.CallbackQuery != nil:
			// Parsing mode: callback
			cnvCtx = parseCallback(tgUpt, user)
		case (tgUpt.Message != nil) && (tgUpt.Message.Text[0] == '/'):
			// Parsing mode: command
			cnvCtx = parseCommand(tgUpt)
		default:
			// Parsing mode: text
			cnvCtx = parseText(tgUpt)
		}

		upt := updates.UserUpdate{
			Interaction: updates.InteractionContext{
				Conversation: coalesceNil(user.Conversation, "root"),
				State:        coalesceNil(user.State, "init"),
				ChatID:       coalesceNil(user.InteractionChatID, 0),
				MessageID:    coalesceNil(user.InteractionMessageID, 0),
			},
			Incoming: cnvCtx,
		}

		return upt, nil
	}
}
