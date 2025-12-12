package domain

import (
	"testing"
	"time"
)

func TestNewUpdatedOn(t *testing.T) {
	t.Parallel()
	before := time.Now().UTC()
	u := NewUpdatedOn()
	after := time.Now().UTC()

	timeVal := time.Time(u)
	if timeVal.Before(before) || timeVal.After(after) {
		t.Errorf("UpdatedOn not in expected range: got %v, want between %v and %v", timeVal, before, after)
	}
}
