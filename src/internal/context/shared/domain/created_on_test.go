package domain

import (
	"testing"
	"time"
)

func TestCreatedOn(t *testing.T) {
	t.Run("NewCreatedOn returns UTC time in range", func(t *testing.T) {
		t.Parallel()
		before := time.Now().UTC()
		c := NewCreatedOn()
		after := time.Now().UTC()

		timeVal := time.Time(c)
		if timeVal.Before(before) || timeVal.After(after) {
			t.Errorf("CreatedOn not in expected range: got %v, want between %v and %v", timeVal, before, after)
		}
	})
}
