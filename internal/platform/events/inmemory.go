package events

import (
	"context"
	"sync"
	"time"
)

type InMemoryBus struct {
	mu          sync.RWMutex
	subscribers map[string][]Handler
}

func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{
		subscribers: make(map[string][]Handler),
	}
}

func (bus *InMemoryBus) Publish(ctx context.Context, event Event) error {
	bus.mu.RLock()
	handlers := append([]Handler(nil), bus.subscribers[event.Name]...)
	bus.mu.RUnlock()

	if event.OccurredAt.IsZero() {
		event.OccurredAt = time.Now().UTC()
	}

	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			return err
		}
	}

	return nil
}

func (bus *InMemoryBus) Subscribe(eventName string, handler Handler) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	bus.subscribers[eventName] = append(bus.subscribers[eventName], handler)
}
