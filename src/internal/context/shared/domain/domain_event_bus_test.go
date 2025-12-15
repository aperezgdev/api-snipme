package domain

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockSubscriber struct {
	called        bool
	receivedEvent DomainEvent
}

func (m *mockSubscriber) On(ctx context.Context, event DomainEvent) error {
	m.called = true
	m.receivedEvent = event
	return nil
}

type testEventWithName struct {
	DomainEventBase
}

func (t testEventWithName) Name() string { return t.DomainEventBase.Name }

func TestEventBusInMemory(t *testing.T) {
	t.Run("Publish calls On for subscribers", func(t *testing.T) {
		eb := NewEventBusInMemory()
		sub := &mockSubscriber{}
		event := testEventWithName{DomainEventBase: NewDomainEvent("aggid", "eventName")}
		eb.AddSubscribers("eventName", sub)
		eb.Publish(context.Background(), event)
		assert.True(t, sub.called)
		assert.Equal(t, event, sub.receivedEvent)
	})

	t.Run("Publish does nothing if no subscribers", func(t *testing.T) {
		eb := NewEventBusInMemory()
		event := testEventWithName{DomainEventBase: NewDomainEvent("aggid", "noSubs")}
		eb.Publish(context.Background(), event)
	})

	t.Run("AddSubscribers appends multiple subscribers", func(t *testing.T) {
		eb := NewEventBusInMemory()
		sub1 := &mockSubscriber{}
		sub2 := &mockSubscriber{}
		eb.AddSubscribers("eventName", sub1, sub2)
		assert.Len(t, eb.subscribers["eventName"], 2)
	})
}
