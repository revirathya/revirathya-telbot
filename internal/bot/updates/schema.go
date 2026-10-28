package updates

type InteractionContext struct {
	Conversation string
	State        string
	ChatID       int64
	MessageID    int
}

type ConversationContext struct {
	ChatID    int64
	MessageID int
	Mode      string
	Key       string
	Data      string
}

type UserUpdate struct {
	Interaction InteractionContext
	Incoming    ConversationContext
}
