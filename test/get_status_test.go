package test

import (
	"net/http"
	"net/url"
	"testing"
)

func TestGetStatusTest(t *testing.T) {
	parsed, err := url.Parse(API_URL)
	if err != nil {
		t.Fatalf("Invalid API_URL: %v", err)
	}

	t.Run("Run perfecly fine", func(t *testing.T) {
		client := http.Client{}
		response, err := client.Do(&http.Request{
			Method: http.MethodGet,
			URL: &url.URL{
				Scheme: parsed.Scheme,
				Host:   parsed.Host,
				Path:   "/status",
			},
		})
		if err != nil {
			t.Fatalf("Error making request: %v", err)
		}
		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
		}
	})
}
