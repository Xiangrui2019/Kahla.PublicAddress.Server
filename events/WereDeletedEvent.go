package events

import "Kahla.PublicAddress.Server/models"

type WereDeletedEvent struct {
	Event
	Trigger models.User `json:"trigger"`
}
