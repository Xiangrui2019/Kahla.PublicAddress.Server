package events

type TimerUpdatedEvent struct {
	Event
	ConversationID int `json:"conversationId"`
	NewTimer       int `json:"newTimer"`
}
