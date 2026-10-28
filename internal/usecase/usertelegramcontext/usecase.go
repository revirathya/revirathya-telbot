package usertelegramcontext

import (
	"context"

	"github.com/revirathya/revirathya-telbot/internal/domain"
)

var (
	rootConversation string = "root"
	rootState        string = "init"
)

// Usecase
type UserTelegramContextUsecase struct {
	userTelegramContext domain.UserTelegramContextRepository
}

func NewUserTelegramContextUsecase(userTelegramContext domain.UserTelegramContextRepository) *UserTelegramContextUsecase {
	return &UserTelegramContextUsecase{
		userTelegramContext: userTelegramContext,
	}
}

func (uc *UserTelegramContextUsecase) GetUserContext(ctx context.Context) (*domain.UserTelegramContext, error) {
	return uc.userTelegramContext.GetUserContext(ctx)
}

func (uc *UserTelegramContextUsecase) StartConversation(ctx context.Context, conversation string, chatID int64, messageID int) (*domain.UserTelegramContext, error) {
	patch := domain.UserTelegramContextPatch{
		Conversation:         &conversation,
		State:                &rootState,
		InteractionChatID:    &chatID,
		InteractionMessageID: &messageID,
	}
	return uc.userTelegramContext.UpdateUserContext(ctx, patch)
}

func (uc *UserTelegramContextUsecase) ChangeState(ctx context.Context, state string) (*domain.UserTelegramContext, error) {
	patch := domain.UserTelegramContextPatch{
		State: &state,
	}
	return uc.userTelegramContext.UpdateUserContext(ctx, patch)
}

func (uc *UserTelegramContextUsecase) CloseConversation(ctx context.Context) (*domain.UserTelegramContext, error) {
	var zeroInt int
	var zeroInt64 int64
	patch := domain.UserTelegramContextPatch{
		Conversation:         &rootConversation,
		State:                &rootState,
		InteractionChatID:    &zeroInt64,
		InteractionMessageID: &zeroInt,
	}
	return uc.userTelegramContext.UpdateUserContext(ctx, patch)
}
