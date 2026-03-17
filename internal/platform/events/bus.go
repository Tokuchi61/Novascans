package events

import (
	"context"
	"time"
)

type Event struct {
	Name       string
	Payload    any
	OccurredAt time.Time
}

type Handler func(context.Context, Event) error

type Bus interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(eventName string, handler Handler)
}
