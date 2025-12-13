package domain

import (
	"testing"
	"time"
)

func TestNewRefreshTokenExpiresAt(t *testing.T) {
	t.Run("Success on future time", func(t *testing.T) {
		t.Parallel()
		futureTime := time.Now().Add(24 * time.Hour)

		expiresAt, err := NewRefreshTokenExpiresAt(futureTime)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if time.Time(expiresAt).IsZero() {
			t.Error("Expected ExpiresAt to be set")
		}
	})

	t.Run("Returns zero value on past time", func(t *testing.T) {
		t.Parallel()
		pastTime := time.Now().Add(-24 * time.Hour)

		expiresAt, err := NewRefreshTokenExpiresAt(pastTime)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !time.Time(expiresAt).IsZero() {
			t.Error("Expected ExpiresAt to be zero for past time")
		}
	})

	t.Run("Returns zero value on current time", func(t *testing.T) {
		t.Parallel()
		currentTime := time.Now().Add(-1 * time.Second)

		expiresAt, err := NewRefreshTokenExpiresAt(currentTime)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !time.Time(expiresAt).IsZero() {
			t.Error("Expected ExpiresAt to be zero for current time")
		}
	})
}
