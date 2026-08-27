package grpc

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"

	financecorev1 "github.com/revirathya/revirathya-telbot/gen/go/financecore/v1"
	"github.com/revirathya/revirathya-telbot/internal/domain"
	grpcCon "github.com/revirathya/revirathya-telbot/internal/platform/grpc"
)

type userTelegramContextRepositoryImpl struct {
	fincoreConn *grpcCon.FinanceCoreConn
}

func NewUserTelegramContextRepositoryImpl(fincoreConn *grpcCon.FinanceCoreConn) domain.UserTelegramContextRepository {
	return &userTelegramContextRepositoryImpl{
		fincoreConn: fincoreConn,
	}
}

func (r *userTelegramContextRepositoryImpl) GetUserContext(ctx context.Context, userID uuid.UUID) (*domain.UserTelegramContext, error) {
	in := &financecorev1.GetUserTelegramContextRequest{
		UserId: userID.String(),
	}

	row, err := r.fincoreConn.UserTelegramContext.GetUserTelegramContext(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed to get user telegram context: %v", err)
	}

	uuidUserID, err := uuid.Parse(row.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid user id UUID format: %v", err)
	}

	var intInteractionChatID *int64
	if row.InteractionChatId != nil {
		cvtInteractionChatID, err := strconv.ParseInt(*row.InteractionChatId, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid interaction chat id Int format: %v", err)
		}

		intInteractionChatID = &cvtInteractionChatID
	}

	var intInteractionMessageID *int
	if row.InteractionMessageId != nil {
		cvtInteractionMessageID, err := strconv.ParseInt(*row.InteractionChatId, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid interaction message id Int format: %v", err)
		}

		intInteractionChatID = &cvtInteractionMessageID
	}

	result := &domain.UserTelegramContext{
		UserID:               uuidUserID,
		Conversation:         row.Conversation,
		State:                row.State,
		InteractionChatID:    intInteractionChatID,
		InteractionMessageID: intInteractionMessageID,
	}
	return result, nil
}

func (r *userTelegramContextRepositoryImpl) GetUserContextByTelegramID(ctx context.Context, telegramChatID string) (*domain.UserTelegramContext, error) {
	in := &financecorev1.GetUserTelegramContextByTelegramIDRequest{
		TelegramChatId: telegramChatID,
	}

	row, err := r.fincoreConn.UserTelegramContext.GetUserTelegramContextByTelegramID(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed to get user telegram context by telegram id: %v", err)
	}

	uuidUserID, err := uuid.Parse(row.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid user id UUID format: %v", err)
	}

	var intInteractionChatID *int64
	if row.InteractionChatId != nil {
		cvtInteractionChatID, err := strconv.ParseInt(*row.InteractionChatId, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid interaction chat id Int format: %v", err)
		}

		intInteractionChatID = &cvtInteractionChatID
	}

	var intInteractionMessageID *int
	if row.InteractionMessageId != nil {
		cvtInteractionMessageID, err := strconv.Atoi(*row.InteractionMessageId)
		if err != nil {
			return nil, fmt.Errorf("invalid interaction message id Int format: %v", err)
		}

		intInteractionMessageID = &cvtInteractionMessageID
	}

	result := &domain.UserTelegramContext{
		UserID:               uuidUserID,
		Conversation:         row.Conversation,
		State:                row.State,
		InteractionChatID:    intInteractionChatID,
		InteractionMessageID: intInteractionMessageID,
	}
	return result, nil
}

func (r *userTelegramContextRepositoryImpl) UpdateUserContext(ctx context.Context, userID uuid.UUID, conversation *string, state *string, interactionChatID *int64, interactionMessageID *int) (*domain.UserTelegramContext, error) {
	var ptrInteractionChatID *string
	if interactionChatID != nil {
		strInteractionChatID := strconv.FormatInt(*interactionChatID, 10)
		ptrInteractionChatID = &strInteractionChatID
	}

	var ptrInteractionMessageID *string
	if interactionMessageID != nil {
		strInteractionMessageID := strconv.Itoa(*interactionMessageID)
		ptrInteractionMessageID = &strInteractionMessageID
	}

	in := &financecorev1.UpdateUserTelegramContextRequest{
		UserId:               userID.String(),
		Conversation:         conversation,
		State:                state,
		InteractionChatId:    ptrInteractionChatID,
		InteractionMessageId: ptrInteractionMessageID,
	}

	row, err := r.fincoreConn.UserTelegramContext.UpdateUserTelegramContext(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed to update user telegram context: %v", err)
	}

	uuidUserID, err := uuid.Parse(row.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid user id UUID format: %v", err)
	}

	var intInteractionChatID *int64
	if row.InteractionChatId != nil {
		cvtInteractionChatID, err := strconv.ParseInt(*row.InteractionChatId, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid interaction chat id Int format: %v", err)
		}

		intInteractionChatID = &cvtInteractionChatID
	}

	var intInteractionMessageID *int
	if row.InteractionMessageId != nil {
		cvtInteractionMessageID, err := strconv.Atoi(*row.InteractionMessageId)
		if err != nil {
			return nil, fmt.Errorf("invalid interaction message id Int format: %v", err)
		}

		intInteractionMessageID = &cvtInteractionMessageID
	}

	result := &domain.UserTelegramContext{
		UserID:               uuidUserID,
		Conversation:         row.Conversation,
		State:                row.State,
		InteractionChatID:    intInteractionChatID,
		InteractionMessageID: intInteractionMessageID,
	}
	return result, nil
}
