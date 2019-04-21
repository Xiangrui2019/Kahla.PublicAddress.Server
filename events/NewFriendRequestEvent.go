package events

import "Kahla.PublicAddress.Server/models"

type NewFriendRequestEvent struct {
	Event
	RequesterID string      `json:"requesterId"`
	Requester   models.User `json:"requester"`
}
