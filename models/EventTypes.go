package models

const (
	EventTypeNewMessage = iota
	EventTypeNewFriendRequest
	EventTypeWereDeleted
	EventTypeFriendAccepted
	EventTypeTimerUpdated
)
