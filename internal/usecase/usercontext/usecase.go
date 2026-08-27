package usercontext

import (
	"context"

	"github.com/google/uuid"

	"github.com/revirathya/revirathya-telbot/internal/domain"
)

// Usecase: User Telegram Context
type UserTelegramContextUsecase struct {
	utelRepo domain.UserTelegramContextRepository
}

func NewUserTelegramContextUsecase(utelRepo domain.UserTelegramContextRepository) *UserTelegramContextUsecase {
	return &UserTelegramContextUsecase{
		utelRepo: utelRepo,
	}
}

func (uc *UserTelegramContextUsecase) GetUserContextByTelegramID(ctx context.Context, telegramChatID string) (*domain.UserTelegramContext, error) {
	return uc.utelRepo.GetUserContextByTelegramID(ctx, telegramChatID)
}

func (uc *UserTelegramContextUsecase) StartConversation(ctx context.Context, userID uuid.UUID, conversation string, chatID int64, messageID int) (*domain.UserTelegramContext, error) {
	state := "init"
	return uc.utelRepo.UpdateUserContext(
		ctx,
		userID,
		&conversation,
		&state,
		&chatID,
		&messageID,
	)
}

func (uc *UserTelegramContextUsecase) ChangeState(ctx context.Context, userID uuid.UUID, state string) (*domain.UserTelegramContext, error) {
	return uc.utelRepo.UpdateUserContext(ctx, userID, nil, &state, nil, nil)
}

func (uc *UserTelegramContextUsecase) CloseConversation(ctx context.Context, userID uuid.UUID) (*domain.UserTelegramContext, error) {
	var defConversation string = "root"
	var defState string = "init"
	var defChatID int64 = 0
	var defMessageID int = 0

	return uc.utelRepo.UpdateUserContext(ctx, userID, &defConversation, &defState, &defChatID, &defMessageID)
}
