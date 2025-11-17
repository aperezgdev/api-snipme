package test

import (
	"net/http"
	"net/url"
	"testing"
)

func TestDeleteShortLinkByCodeTest(t *testing.T) {
	parsed, err := url.Parse(API_URL)
	if err != nil {
		t.Fatalf("Invalid API_URL: %v", err)
	}
	client := http.Client{}

	t.Run("Success - valid code", func(t *testing.T) {
		response, err := client.Do(&http.Request{
			Method: http.MethodDelete,
			URL: &url.URL{
				Scheme: parsed.Scheme,
				Host:   parsed.Host,
				Path:   "/short-link/66666666-6666-6666-6666-666666666666",
			},
		})
		if err != nil {
			t.Fatalf("Error making request: %v", err)
		}
		if response.StatusCode != http.StatusNoContent {
			t.Fatalf("Expected status code %d, got %d", http.StatusNoContent, response.StatusCode)
		}
	})

	t.Run("Error - not existing code", func(t *testing.T) {
		response, err := client.Do(&http.Request{
			Method: http.MethodDelete,
			URL: &url.URL{
				Scheme: parsed.Scheme,
				Host:   parsed.Host,
				Path:   "/short-link/66666666-6666-6666-6666-666666666656",
			},
		})
		if err != nil {
			t.Fatalf("Error making request: %v", err)
		}
		if response.StatusCode != http.StatusNotFound {
			t.Fatalf("Expected status code %d, got %d", http.StatusNotFound, response.StatusCode)
		}
	})
}
