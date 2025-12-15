package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDomainEventBase(t *testing.T) {
	id := "id"
	aggId := "aggid"
	name := "event"
	occurred := time.Now().UTC()
	base := DomainEventBase{Id: id, AggregateId: aggId, Name: name, OcurredOn: occurred}

	t.Run("ID returns Id", func(t *testing.T) {
		assert.Equal(t, id, base.ID())
	})
	t.Run("AggregateID returns AggregateId", func(t *testing.T) {
		assert.Equal(t, aggId, base.AggregateID())
	})
	t.Run("OccurredOn returns OcurredOn", func(t *testing.T) {
		assert.Equal(t, occurred, base.OccurredOn())
	})
}

func TestNewDomainEvent(t *testing.T) {
	t.Run("returns DomainEventBase with correct fields", func(t *testing.T) {
		aggregateId := "aggid"
		name := "event"
		de := NewDomainEvent(aggregateId, name)
		assert.Equal(t, aggregateId, de.AggregateId)
		assert.Equal(t, name, de.Name)
		assert.NotEmpty(t, de.Id)
		assert.WithinDuration(t, time.Now().UTC(), de.OcurredOn, time.Second)
	})
}
