package domain

// Model: UserContext
type UserContext struct {
	Conversation         string
	State                string
	InteractionChatID    int64
	InteractionMessageID int
}

func NewUserContext() *UserContext {
	return &UserContext{
		Conversation:         "root",
		State:                "init",
		InteractionChatID:    0,
		InteractionMessageID: 0,
	}
}

func (uc *UserContext) StartConversation(conversation string, interactionChatID int64, interactionMessageID int) {
	uc.Conversation = conversation
	uc.State = "init"
	uc.InteractionChatID = interactionChatID
	uc.InteractionMessageID = interactionMessageID
}

func (uc *UserContext) ChangeState(state string) {
	uc.State = state
}

func (uc *UserContext) CloseConversation() {
	uc.Conversation = "root"
	uc.State = "init"
	uc.InteractionChatID = 0
	uc.InteractionMessageID = 0
}
