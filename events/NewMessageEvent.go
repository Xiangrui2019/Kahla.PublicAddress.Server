package events

import "Kahla.PublicAddress.Server/models"

type NewMessageEvent struct {
	Event
	ConversationID int         `json:"conversationId"`
	Sender         models.User `json:"sender"`
	Content        string      `json:"content"`
	AesKey         string      `json:"aesKey"`
	Muted          bool        `json:"muted"`
	SentByMe       bool        `json:"sentByMe"`
}
