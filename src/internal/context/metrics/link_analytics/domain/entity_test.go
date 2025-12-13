package domain

import (
	"testing"
	"time"
)

func TestNewLinkAnalytics(t *testing.T) {
	t.Run("Success on valid linkID", func(t *testing.T) {
		t.Parallel()
		linkID := "00000000-0000-0000-0000-000000000000"

		linkAnalytics, err := NewLinkAnalytics(linkID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if linkAnalytics.LinkId.String() != linkID {
			t.Errorf("Expected LinkId %s, got %s", linkID, linkAnalytics.LinkId.String())
		}
		if linkAnalytics.Id.String() == "" {
			t.Error("Expected ID to be generated")
		}
		if uint(linkAnalytics.TotalViews) != 0 {
			t.Errorf("Expected TotalViews 0, got %d", uint(linkAnalytics.TotalViews))
		}
		if uint(linkAnalytics.UniqueViews) != 0 {
			t.Errorf("Expected UniqueViews 0, got %d", uint(linkAnalytics.UniqueViews))
		}
		if len(linkAnalytics.LinkCountriesViewCounters) != 0 {
			t.Errorf("Expected empty LinkCountriesViewCounters, got %d items", len(linkAnalytics.LinkCountriesViewCounters))
		}
		if time.Time(linkAnalytics.UpdateOn).IsZero() {
			t.Error("Expected UpdateOn to be set")
		}
	})

	t.Run("Fails on invalid linkID", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name    string
			linkID  string
		}{
			{
				name:   "invalid UUID format",
				linkID: "invalid-uuid",
			},
			{
				name:   "empty string",
				linkID: "",
			},
			{
				name:   "malformed UUID",
				linkID: "00000000-0000-0000-0000",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				_, err := NewLinkAnalytics(tt.linkID)
				if err == nil {
					t.Fatalf("Expected error for %s, got nil", tt.name)
				}
			})
		}
	})
}

func TestLinkAnalytics_IncrementTotalViews(t *testing.T) {
	t.Run("Increments total views correctly", func(t *testing.T) {
		t.Parallel()
		linkID := "00000000-0000-0000-0000-000000000000"
		linkAnalytics, _ := NewLinkAnalytics(linkID)

		initialViews := uint(linkAnalytics.TotalViews)
		linkAnalytics.IncrementTotalViews(5)

		if uint(linkAnalytics.TotalViews) != initialViews+5 {
			t.Errorf("Expected TotalViews %d, got %d", initialViews+5, uint(linkAnalytics.TotalViews))
		}
	})

	t.Run("Increments multiple times", func(t *testing.T) {
		t.Parallel()
		linkID := "00000000-0000-0000-0000-000000000000"
		linkAnalytics, _ := NewLinkAnalytics(linkID)

		linkAnalytics.IncrementTotalViews(3)
		linkAnalytics.IncrementTotalViews(7)
		linkAnalytics.IncrementTotalViews(2)

		if uint(linkAnalytics.TotalViews) != 12 {
			t.Errorf("Expected TotalViews 12, got %d", uint(linkAnalytics.TotalViews))
		}
	})

	t.Run("Increments by zero", func(t *testing.T) {
		t.Parallel()
		linkID := "00000000-0000-0000-0000-000000000000"
		linkAnalytics, _ := NewLinkAnalytics(linkID)

		initialViews := uint(linkAnalytics.TotalViews)
		linkAnalytics.IncrementTotalViews(0)

		if uint(linkAnalytics.TotalViews) != initialViews {
			t.Errorf("Expected TotalViews %d, got %d", initialViews, uint(linkAnalytics.TotalViews))
		}
	})
}

func TestLinkAnalytics_IncrementUniqueViews(t *testing.T) {
	t.Run("Increments unique views correctly", func(t *testing.T) {
		t.Parallel()
		linkID := "00000000-0000-0000-0000-000000000000"
		linkAnalytics, _ := NewLinkAnalytics(linkID)

		initialViews := uint(linkAnalytics.UniqueViews)
		linkAnalytics.IncrementUniqueViews(3)

		if uint(linkAnalytics.UniqueViews) != initialViews+3 {
			t.Errorf("Expected UniqueViews %d, got %d", initialViews+3, uint(linkAnalytics.UniqueViews))
		}
	})

	t.Run("Increments multiple times", func(t *testing.T) {
		t.Parallel()
		linkID := "00000000-0000-0000-0000-000000000000"
		linkAnalytics, _ := NewLinkAnalytics(linkID)

		linkAnalytics.IncrementUniqueViews(2)
		linkAnalytics.IncrementUniqueViews(5)
		linkAnalytics.IncrementUniqueViews(1)

		if uint(linkAnalytics.UniqueViews) != 8 {
			t.Errorf("Expected UniqueViews 8, got %d", uint(linkAnalytics.UniqueViews))
		}
	})

	t.Run("Increments by zero", func(t *testing.T) {
		t.Parallel()
		linkID := "00000000-0000-0000-0000-000000000000"
		linkAnalytics, _ := NewLinkAnalytics(linkID)

		initialViews := uint(linkAnalytics.UniqueViews)
		linkAnalytics.IncrementUniqueViews(0)

		if uint(linkAnalytics.UniqueViews) != initialViews {
			t.Errorf("Expected UniqueViews %d, got %d", initialViews, uint(linkAnalytics.UniqueViews))
		}
	})
}
