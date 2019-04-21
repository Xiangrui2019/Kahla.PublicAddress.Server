package events

import "Kahla.PublicAddress.Server/models"

type FriendAcceptedEvent struct {
	Event
	Target models.User `json:"target"`
}
