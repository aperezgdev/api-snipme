package domain

type AggregateRoot struct {
	domainEvents []DomainEvent
}

func (ar *AggregateRoot) Record(event DomainEvent) {
	if ar.domainEvents == nil {
		ar.domainEvents = []DomainEvent{event}
		return
	}
	ar.domainEvents = append(ar.domainEvents, event)
}

func (ar *AggregateRoot) PullDomainEvents() []DomainEvent {
	domainEvents := ar.domainEvents
	ar.domainEvents = []DomainEvent{}
	return domainEvents
}
