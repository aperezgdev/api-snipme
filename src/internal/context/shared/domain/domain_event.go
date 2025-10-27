package domain

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type DomainEventBase struct {
	Id          string
	AggregateId string
	Name        string
	OcurredOn   time.Time
}

func (de DomainEventBase) ID() string {
	return de.Id
}

func (de DomainEventBase) AggregateID() string {
	return de.AggregateId
}

func (de DomainEventBase) OccurredOn() time.Time {
	return de.OcurredOn
}

type DomainEvent interface {
	ID() string
	Name() string
	AggregateID() string
	OccurredOn() time.Time
}

func NewDomainEvent(aggregateId string, name string) DomainEventBase {
	idVO, _ := NewID()
	return DomainEventBase{
		Id:          idVO.String(),
		AggregateId: aggregateId,
		Name:        name,
		OcurredOn:   time.Now().UTC(),
	}
}

type DomainEventSubscriber interface {
	On(context.Context, DomainEvent)
}

type EventBus interface {
	Publish(context.Context, ...DomainEvent)
	AddSubscribers(string, ...DomainEventSubscriber)
}

type EventBusInMemory struct {
	subscribers map[string][]DomainEventSubscriber
}

func NewEventBusInMemory() EventBusInMemory {
	return EventBusInMemory{
		subscribers: make(map[string][]DomainEventSubscriber),
	}
}

func (eb *EventBusInMemory) Publish(ctx context.Context, events ...DomainEvent) {
	for _, event := range events {
		if subs, ok := eb.subscribers[event.Name()]; ok {
			for _, sub := range subs {
				sub.On(ctx, event)
			}
		}
	}
}

func (eb *EventBusInMemory) AddSubscribers(eventName string, subscribers ...DomainEventSubscriber) {
	eb.subscribers[eventName] = append(eb.subscribers[eventName], subscribers...)
}

type EventBusMock struct {
	mock.Mock
}

func (eb *EventBusMock) Publish(ctx context.Context, events ...DomainEvent) {
	eb.Called(ctx, events)
}

func (eb *EventBusMock) AddSubscribers(eventName string, subscribers ...DomainEventSubscriber) {
	eb.Called(eventName, subscribers)
}
