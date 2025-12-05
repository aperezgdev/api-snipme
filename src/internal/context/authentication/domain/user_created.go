package domain

import (
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

const UserCreatedEventName = "UserCreated"

type UserCreated struct {
	shared_domain.DomainEventBase
	name  string
	email string
}

func NewUserCreatedEvent(userId, email, name string) shared_domain.DomainEvent {
	return UserCreated{
		DomainEventBase: shared_domain.NewDomainEvent(
			userId,
			UserCreatedEventName,
		),
		email: email,
		name:  name,
	}
}

func (e UserCreated) Name() string {
	return UserCreatedEventName
}

func (e UserCreated) Email() string {
	return e.email
}

func (e UserCreated) UserName() string {
	return e.name
}
