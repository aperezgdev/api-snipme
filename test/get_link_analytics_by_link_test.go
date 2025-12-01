package test

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
)

type linkAnalyticsResponse struct {
    LinkId      string `json:"linkId"`
    TotalViews  uint   `json:"totalViews"`
    UniqueViews uint   `json:"uniqueViews"`
}

func TestGetLinkAnalyticsByLinkTest(t *testing.T) {
    client := http.Client{}
    parsed, err := url.Parse(API_URL)
    if err != nil {
        t.Fatalf("Invalid API_URL: %v", err)
    }

    t.Run("Success - valid link_id", func(t *testing.T) {
        req := &http.Request{
            Method: http.MethodGet,
            URL: &url.URL{
                Scheme:   parsed.Scheme,
                Host:     parsed.Host,
                Path:     "/link-analytics",
                RawQuery: url.Values{
                    "link_id": []string{"22222222-2222-2222-2222-222222222222"},
                }.Encode(),
            },
        }

        resp, err := client.Do(req)
        if err != nil {
            t.Fatalf("Error making request: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
        }

        var response linkAnalyticsResponse
        err = json.NewDecoder(resp.Body).Decode(&response)
        if err != nil {
            t.Fatalf("Error decoding response: %v", err)
        }

        if response.LinkId == "" {
            t.Fatal("Expected linkId to be set")
        }
    })

    t.Run("Error - link_id not found", func(t *testing.T) {
        req := &http.Request{
            Method: http.MethodGet,
            URL: &url.URL{
                Scheme:   parsed.Scheme,
                Host:     parsed.Host,
                Path:     "/link-analytics",
                RawQuery: url.Values{
                    "link_id": []string{"33333333-3333-3333-3333-333333333333"},
                }.Encode(),
            },
        }

        resp, err := client.Do(req)
        if err != nil {
            t.Fatalf("Error making request: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusNotFound {
            t.Fatalf("Expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
        }
    })

    t.Run("Error - missing link_id parameter", func(t *testing.T) {
        req := &http.Request{
            Method: http.MethodGet,
            URL: &url.URL{
                Scheme: parsed.Scheme,
                Host:   parsed.Host,
                Path:   "/link-analytics",
            },
        }

        resp, err := client.Do(req)
        if err != nil {
            t.Fatalf("Error making request: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusBadRequest {
            t.Fatalf("Expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
        }
    })

    t.Run("Error - invalid link_id format", func(t *testing.T) {
        req := &http.Request{
            Method: http.MethodGet,
            URL: &url.URL{
                Scheme:   parsed.Scheme,
                Host:     parsed.Host,
                Path:     "/link-analytics",
                RawQuery: "link_id=invalid-uuid",
            },
        }

        resp, err := client.Do(req)
        if err != nil {
            t.Fatalf("Error making request: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusBadRequest {
            t.Fatalf("Expected status code %d, got %d",http.StatusBadRequest, resp.StatusCode)
        }
    })
}