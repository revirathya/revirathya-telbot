package usertelegramcontext

import (
	"context"

	financecorev1 "github.com/revirathya/revirathya-finance-core/gen/grpc/financecore/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// Internal
	"github.com/revirathya/revirathya-telbot/internal/domain"
	errDom "github.com/revirathya/revirathya-telbot/internal/domain/error"
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

// Private
func (r *userTelegramContextRepositoryImpl) protoToDomain(usr *financecorev1.UserTelegramContext) (*domain.UserTelegramContext, error) {
	interactionChatID, err := converter.FromOptionalStringToInt64("interaction_chat_id", usr.InteractionChatId)
	if err != nil {
		return nil, err
	}

	interactionMessageID, err := converter.FromOptionalStringToInt("interaction_message_id", usr.InteractionMessageId)
	if err != nil {
		return nil, err
	}

	result := &domain.UserTelegramContext{
		Conversation:         usr.Conversation,
		State:                usr.State,
		InteractionChatID:    interactionChatID,
		InteractionMessageID: interactionMessageID,
	}
	return result, nil
}

// Public
func (r *userTelegramContextRepositoryImpl) GetUserContext(ctx context.Context) (*domain.UserTelegramContext, error) {
	// Prepare input
	in := &financecorev1.GetUserTelegramContextRequest{}

	// Get user telegram context
	row, err := r.finance.UserTelegramContext.GetUserTelegramContext(ctx, in)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, errDom.ErrNotFound
		}
		return nil, err
	}

	// Parse result
	return r.protoToDomain(row.UserTelegramContext)
}

func (r *userTelegramContextRepositoryImpl) UpdateUserContext(ctx context.Context, patch domain.UserTelegramContextPatch) (*domain.UserTelegramContext, error) {
	// Prepare input
	in := &financecorev1.UpdateUserTelegramContextRequest{
		Conversation:         patch.Conversation,
		State:                patch.State,
		InteractionChatId:    converter.FromOptionalInt64ToString("interaction_chat_id", patch.InteractionChatID),
		InteractionMessageId: converter.FromOptionalIntToString("interaction_message_id", patch.InteractionMessageID),
	}

	// Update user telegram context
	_, err := r.finance.UserTelegramContext.UpdateUserTelegramContext(ctx, in)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, errDom.ErrNotFound
		}
		if status.Code(err) == codes.InvalidArgument {
			return nil, errDom.ErrInvalidInput
		}
		return nil, err
	}

	return r.GetUserContext(ctx)
}
