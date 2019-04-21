package models

import "fmt"

type InvalidEventTypeError struct {
	EventType int
}

func (i *InvalidEventTypeError) Error() string {
	return fmt.Sprintf("invalid event type: %d", i.EventType)
}