package test

import (
	"net/http"
	"net/url"
	"testing"
)

func TestGetShortLinkByCodeTest(t *testing.T) {
	parsed, err := url.Parse(API_URL)
	if err != nil {
		t.Fatalf("Invalid API_URL: %v", err)
	}
	client := http.Client{}

	t.Run("Success - valid code", func(t *testing.T) {
		response, err := client.Do(&http.Request{
			Method: http.MethodGet,
			URL: &url.URL{
				Scheme: parsed.Scheme,
				Host:   parsed.Host,
				Path:   "/testcode",
			},
		})
		if err != nil {
			t.Fatalf("Error making request: %v", err)
		}
		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
		}
	})

	t.Run("Error - not existing code", func(t *testing.T) {
		response, err := client.Do(&http.Request{
			Method: http.MethodGet,
			URL: &url.URL{
				Scheme: parsed.Scheme,
				Host:   parsed.Host,
				Path:   "/notexist",
			},
		})
		if err != nil {
			t.Fatalf("Error making request: %v", err)
		}
		if response.StatusCode != http.StatusNotFound {
			t.Fatalf("Expected status code %d, got %d", http.StatusBadRequest, response.StatusCode)
		}
	})
}
