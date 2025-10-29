package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"
)

type postShortLinkRequest struct {
	OriginalURL string `json:"original_url"`
	ClientID    string `json:"client_id"`
}

func TestPostShortLinkTest(t *testing.T) {
	client := http.Client{}
	parsed, err := url.Parse(API_URL)
	if err != nil {
		t.Fatalf("Invalid API_URL: %v", err)
	}

	t.Run("Success - create short link", func(t *testing.T) {
		body := postShortLinkRequest{
			OriginalURL: "https://example.com",
			ClientID:    "11111111-1111-1111-1111-111111111111",
		}
		b, _ := json.Marshal(body)
		req := &http.Request{
			Method: http.MethodPost,
			URL: &url.URL{
				Scheme: parsed.Scheme,
				Host:   parsed.Host,
				Path:   "/short-links",
			},
			Body:          io.NopCloser(bytes.NewReader(b)),
			ContentLength: int64(len(b)),
			Header:        make(http.Header),
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Error making request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 201 or 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Error - invalid body", func(t *testing.T) {
		b := []byte(`{"original_url": "", "client_id": ""}`)
		req := &http.Request{
			Method: http.MethodPost,
			URL: &url.URL{
				Scheme: parsed.Scheme,
				Host:   parsed.Host,
				Path:   "/short-links",
			},
			Body:          io.NopCloser(bytes.NewReader(b)),
			ContentLength: int64(len(b)),
			Header:        make(http.Header),
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Error making request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("Expected status code 400, got %d", resp.StatusCode)
		}
	})
}
