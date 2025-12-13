package domain

import (
	"testing"

	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

func TestNewLinkCountryViewCounter(t *testing.T) {
	t.Run("Success on valid parameters", func(t *testing.T) {
		t.Parallel()
		linkId := "00000000-0000-0000-0000-000000000000"
		countryCode := "US"

		counter, err := NewLinkCountryViewCounter(linkId, countryCode)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if counter.LinkId.String() != linkId {
			t.Errorf("Expected LinkId %s, got %s", linkId, counter.LinkId.String())
		}
		if string(counter.CountryCode) != countryCode {
			t.Errorf("Expected CountryCode %s, got %s", countryCode, string(counter.CountryCode))
		}
		if counter.Id.String() == "" {
			t.Error("Expected ID to be generated")
		}
		if uint(counter.TotalViews) != 0 {
			t.Errorf("Expected TotalViews 0, got %d", uint(counter.TotalViews))
		}
		if uint(counter.UniqueViews) != 0 {
			t.Errorf("Expected UniqueViews 0, got %d", uint(counter.UniqueViews))
		}
	})

	t.Run("Success on different country codes", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name        string
			countryCode string
		}{
			{
				name:        "US code",
				countryCode: "US",
			},
			{
				name:        "GB code",
				countryCode: "GB",
			},
			{
				name:        "ES code",
				countryCode: "ES",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				linkId := "00000000-0000-0000-0000-000000000000"
				counter, err := NewLinkCountryViewCounter(linkId, tt.countryCode)
				if err != nil {
					t.Fatalf("Expected no error for %s, got %v", tt.name, err)
				}
				if string(counter.CountryCode) != tt.countryCode {
					t.Errorf("Expected CountryCode %s, got %s", tt.countryCode, string(counter.CountryCode))
				}
			})
		}
	})

	t.Run("Fails on invalid linkId", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name   string
			linkId string
		}{
			{
				name:   "invalid UUID format",
				linkId: "invalid-uuid",
			},
			{
				name:   "empty string",
				linkId: "",
			},
			{
				name:   "malformed UUID",
				linkId: "00000000-0000-0000-0000",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				_, err := NewLinkCountryViewCounter(tt.linkId, "US")
				if err == nil {
					t.Fatalf("Expected error for %s, got nil", tt.name)
				}
			})
		}
	})

	t.Run("Fails on invalid country code", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name        string
			countryCode string
			expected    string
		}{
			{
				name:        "empty code",
				countryCode: "",
				expected:    "Country code must be 2 characters long",
			},
			{
				name:        "single character",
				countryCode: "U",
				expected:    "Country code must be 2 characters long",
			},
			{
				name:        "three characters",
				countryCode: "USA",
				expected:    "Country code must be 2 characters long",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				linkId := "00000000-0000-0000-0000-000000000000"
				_, err := NewLinkCountryViewCounter(linkId, tt.countryCode)
				if err == nil {
					t.Fatalf("Expected error for %s, got nil", tt.name)
				}
				validationErr, ok := err.(shared_domain.ValidationError)
				if !ok {
					t.Fatalf("Expected ValidationError, got %T", err)
				}
				if validationErr.Message != tt.expected {
					t.Errorf("Expected error message '%s', got '%s'", tt.expected, validationErr.Message)
				}
			})
		}
	})
}

func TestLinkCountryViewCounter_Increment(t *testing.T) {
	t.Run("Increments views correctly", func(t *testing.T) {
		t.Parallel()
		linkId := "00000000-0000-0000-0000-000000000000"
		counter, _ := NewLinkCountryViewCounter(linkId, "US")

		result := counter.Increment(5, 3)

		if uint(result.TotalViews) != 5 {
			t.Errorf("Expected TotalViews 5, got %d", uint(result.TotalViews))
		}
		if uint(result.UniqueViews) != 3 {
			t.Errorf("Expected UniqueViews 3, got %d", uint(result.UniqueViews))
		}
	})

	t.Run("Increments multiple times", func(t *testing.T) {
		t.Parallel()
		linkId := "00000000-0000-0000-0000-000000000000"
		counter, _ := NewLinkCountryViewCounter(linkId, "US")

		counter.Increment(3, 2)
		counter.Increment(7, 4)
		result := counter.Increment(5, 1)

		if uint(result.TotalViews) != 15 {
			t.Errorf("Expected TotalViews 15, got %d", uint(result.TotalViews))
		}
		if uint(result.UniqueViews) != 7 {
			t.Errorf("Expected UniqueViews 7, got %d", uint(result.UniqueViews))
		}
	})

	t.Run("Increments by zero", func(t *testing.T) {
		t.Parallel()
		linkId := "00000000-0000-0000-0000-000000000000"
		counter, _ := NewLinkCountryViewCounter(linkId, "US")

		result := counter.Increment(0, 0)

		if uint(result.TotalViews) != 0 {
			t.Errorf("Expected TotalViews 0, got %d", uint(result.TotalViews))
		}
		if uint(result.UniqueViews) != 0 {
			t.Errorf("Expected UniqueViews 0, got %d", uint(result.UniqueViews))
		}
	})

	t.Run("Returns self for chaining", func(t *testing.T) {
		t.Parallel()
		linkId := "00000000-0000-0000-0000-000000000000"
		counter, _ := NewLinkCountryViewCounter(linkId, "US")

		result := counter.Increment(5, 3)

		if result != counter {
			t.Error("Expected Increment to return the same instance")
		}
	})
}
