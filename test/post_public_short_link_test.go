package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"
)

type postPublicShortLinkRequest struct {
	OriginalLink string `json:"original_link"`
}

func TestPostPublicShortLink(t *testing.T) {
	client := http.Client{}
	parsed, err := url.Parse(API_URL)
	if err != nil {
		t.Fatalf("Invalid API_URL: %v", err)
	}

	t.Run("Run success on valid original link", func(t *testing.T) {
		t.Parallel()
		var postPublicShortLinkRequestBody = postPublicShortLinkRequest{
			OriginalLink: "https://example.com/public",
		}
		b, _ := json.Marshal(postPublicShortLinkRequestBody)
		req := &http.Request{
			Method: http.MethodPost,
			URL: &url.URL{
				Scheme: parsed.Scheme,
				Host:   parsed.Host,
				Path:   "/public-short-links",
			},
			Body:          io.NopCloser(bytes.NewReader(b)),
			ContentLength: int64(len(b)),
			Header:        make(http.Header),
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+TEST_TOKEN)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Error making request: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 201 or 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Run error on invalid original link", func(t *testing.T) {
		t.Parallel()
		var postPublicShortLinkRequestBody = postPublicShortLinkRequest{
			OriginalLink: "invalid-url",
		}
		b, _ := json.Marshal(postPublicShortLinkRequestBody)
		req := &http.Request{
			Method: http.MethodPost,
			URL: &url.URL{
				Scheme: parsed.Scheme,
				Host:   parsed.Host,
				Path:   "/public-short-links",
			},
			Body:          io.NopCloser(bytes.NewReader(b)),
			ContentLength: int64(len(b)),
			Header:        make(http.Header),
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+TEST_TOKEN)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Error making request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("Expected status code 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Run error on empty original link", func(t *testing.T) {
		t.Parallel()
		var postPublicShortLinkRequestBody = postPublicShortLinkRequest{
			OriginalLink: "",
		}
		b, _ := json.Marshal(postPublicShortLinkRequestBody)
		req := &http.Request{
			Method: http.MethodPost,
			URL: &url.URL{
				Scheme: parsed.Scheme,
				Host:   parsed.Host,
				Path:   "/public-short-links",
			},
			Body:          io.NopCloser(bytes.NewReader(b)),
			ContentLength: int64(len(b)),
			Header:        make(http.Header),
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+TEST_TOKEN)

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