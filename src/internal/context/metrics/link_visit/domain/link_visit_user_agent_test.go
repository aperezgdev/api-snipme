package domain

import (
	"testing"
)

func TestNewLinkVisitUserAgent(t *testing.T) {
	t.Run("Creates user agent correctly", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name      string
			userAgent string
		}{
			{
				name:      "Mozilla Firefox",
				userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0",
			},
			{
				name:      "Chrome",
				userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
			},
			{
				name:      "Empty user agent",
				userAgent: "",
			},
			{
				name:      "Simple string",
				userAgent: "MyBot/1.0",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				userAgent := NewLinkVisitUserAgent(tt.userAgent)
				if string(userAgent) != tt.userAgent {
					t.Errorf("Expected user agent %s, got %s", tt.userAgent, string(userAgent))
				}
			})
		}
	})
}
