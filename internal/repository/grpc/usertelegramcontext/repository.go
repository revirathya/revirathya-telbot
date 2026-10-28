package usertelegramcontext

import (
	"context"

	financecorev1 "github.com/revirathya/revirathya-finance-core/gen/grpc/financecore/v1"

	// Internal
	"github.com/revirathya/revirathya-telbot/internal/domain"
	"github.com/revirathya/revirathya-telbot/internal/platform/grpc/clients"
	"github.com/revirathya/revirathya-telbot/internal/platform/grpc/converter"
)

type userTelegramContextRepositoryImpl struct {
	finance *clients.FinanceClient
}

func NewUserTelegramContextRepositoryImpl(finance *clients.FinanceClient) domain.UserTelegramContextRepository {
	return &userTelegramContextRepositoryImpl{
		finance: finance,
	}
}

func (r *userTelegramContextRepositoryImpl) GetUserContext(ctx context.Context) (*domain.UserTelegramContext, error) {
	in := &financecorev1.GetUserTelegramContextRequest{}

	row, err := r.finance.UserTelegramContext.GetUserTelegramContext(ctx, in)
	if err != nil {
		return nil, err
	}

	interactionChatID, err := converter.FromOptionalStringToInt64("interaction_chat_id", row.UserTelegramContext.InteractionChatId)
	if err != nil {
		return nil, err
	}

	interactionMessageID, err := converter.FromOptionalStringToInt("interaction_message_id", row.UserTelegramContext.InteractionMessageId)
	if err != nil {
		return nil, err
	}

	result := &domain.UserTelegramContext{
		Conversation:         row.UserTelegramContext.Conversation,
		State:                row.UserTelegramContext.State,
		InteractionChatID:    interactionChatID,
		InteractionMessageID: interactionMessageID,
	}
	return result, nil
}

func (r *userTelegramContextRepositoryImpl) UpdateUserContext(ctx context.Context, patch domain.UserTelegramContextPatch) (*domain.UserTelegramContext, error) {
	in := &financecorev1.UpdateUserTelegramContextRequest{
		Conversation:         patch.Conversation,
		State:                patch.State,
		InteractionChatId:    converter.FromOptionalInt64ToString("interaction_chat_id", patch.InteractionChatID),
		InteractionMessageId: converter.FromOptionalIntToString("interaction_message_id", patch.InteractionMessageID),
	}

	_, err := r.finance.UserTelegramContext.UpdateUserTelegramContext(ctx, in)
	if err != nil {
		return nil, err
	}

	return r.GetUserContext(ctx)
}
