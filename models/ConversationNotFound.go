package models

type ConversationNotFound struct{}

func (*ConversationNotFound) Error() string {
	return "conversation not found"
}
