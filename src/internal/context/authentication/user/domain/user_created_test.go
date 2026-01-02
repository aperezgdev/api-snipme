package domain

import (
	"testing"
)

func TestNewUserCreatedEvent(t *testing.T) {
	t.Run("Creates event correctly", func(t *testing.T) {
		t.Parallel()
		userId := "00000000-0000-0000-0000-000000000001"
		email := "user@example.com"
		name := "google"

		event := NewUserCreatedEvent(userId, email, name)

		if event.AggregateID() != userId {
			t.Errorf("Expected AggregateID %s, got %s", userId, event.AggregateID())
		}
		if event.Name() != UserCreatedEventName {
			t.Errorf("Expected Name %s, got %s", UserCreatedEventName, event.Name())
		}
		if event.OccurredOn().IsZero() {
			t.Error("Expected OccurredOn to be set")
		}
	})

	t.Run("UserCreated has correct methods", func(t *testing.T) {
		t.Parallel()
		userId := "00000000-0000-0000-0000-000000000001"
		email := "user@example.com"
		name := "google"

		event := NewUserCreatedEvent(userId, email, name)
		userCreated, ok := event.(UserCreated)
		if !ok {
			t.Fatal("Expected UserCreated type")
		}

		if userCreated.Email() != email {
			t.Errorf("Expected Email %s, got %s", email, userCreated.Email())
		}
		if userCreated.UserName() != name {
			t.Errorf("Expected UserName %s, got %s", name, userCreated.UserName())
		}
	})
}

func TestUserCreated_Name(t *testing.T) {
	t.Run("Returns correct event name", func(t *testing.T) {
		t.Parallel()
		event := UserCreated{}
		if event.Name() != UserCreatedEventName {
			t.Errorf("Expected Name %s, got %s", UserCreatedEventName, event.Name())
		}
		if event.Name() != "UserCreated" {
			t.Errorf("Expected Name 'UserCreated', got %s", event.Name())
		}
	})
}
