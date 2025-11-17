package test

import (
	"net/http"
	"net/url"
	"testing"
)

func TestGetShortLinkByClient(t *testing.T) {
	parsed, err := url.Parse(API_URL)
	if err != nil {
		t.Fatalf("Invalid API_URL: %v", err)
	}
	client := http.Client{}

	t.Run("Success - valid client ID", func(t *testing.T) {
		response, err := client.Do(&http.Request{
			Method: http.MethodGet,
			URL: &url.URL{
				Scheme: parsed.Scheme,
				Host:   parsed.Host,
				Path:   "/short-link",
				RawQuery: url.Values{
					"client_id": []string{"55555555-5555-5555-5555-555555555555"},
				}.Encode(),
			},
		})
		if err != nil {
			t.Fatalf("Error making request: %v", err)
		}
		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
		}

	})

	t.Run("Error - invalid client ID", func(t *testing.T) {
		response, err := client.Do(&http.Request{
			Method: http.MethodGet,
			URL: &url.URL{
				Scheme: parsed.Scheme,
				Host:   parsed.Host,
				Path:   "/short-link",
				RawQuery: url.Values{
					"client_id": []string{"44444444-4444-4444-4444-444444444444"},
				}.Encode(),
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
