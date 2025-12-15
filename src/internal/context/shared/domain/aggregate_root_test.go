package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testEvent struct{}

func (testEvent) ID() string            { return "id" }
func (testEvent) Name() string          { return "test" }
func (testEvent) AggregateID() string   { return "aggid" }
func (testEvent) OccurredOn() time.Time { return time.Now().UTC() }

func TestAggregateRoot(t *testing.T) {
	t.Run("Record adds event to domainEvents", func(t *testing.T) {
		ar := &AggregateRoot{}
		event := testEvent{}
		ar.Record(event)
		assert.Len(t, ar.domainEvents, 1)
		assert.Equal(t, event, ar.domainEvents[0])
	})

	t.Run("Record appends multiple events", func(t *testing.T) {
		ar := &AggregateRoot{}
		e1 := testEvent{}
		e2 := testEvent{}
		ar.Record(e1)
		ar.Record(e2)
		assert.Len(t, ar.domainEvents, 2)
		assert.Equal(t, e1, ar.domainEvents[0])
		assert.Equal(t, e2, ar.domainEvents[1])
	})

	t.Run("PullDomainEvents returns and clears events", func(t *testing.T) {
		ar := &AggregateRoot{}
		e1 := testEvent{}
		e2 := testEvent{}
		ar.Record(e1)
		ar.Record(e2)
		events := ar.PullDomainEvents()
		assert.Len(t, events, 2)
		assert.Len(t, ar.domainEvents, 0)
	})
}
