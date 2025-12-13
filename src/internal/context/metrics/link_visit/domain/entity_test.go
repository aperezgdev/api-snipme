package domain

import (
	"testing"
	"time"
)

func TestNewLinkVisit(t *testing.T) {
	t.Run("Success on valid parameters", func(t *testing.T) {
		t.Parallel()
		linkID := "00000000-0000-0000-0000-000000000001"
		visitorId := "00000000-0000-0000-0000-000000000002"
		ip := "192.168.1.1"
		userAgent := "Mozilla/5.0"

		linkVisit, err := NewLinkVisit(linkID, visitorId, ip, userAgent)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if linkVisit.LinkId.String() != linkID {
			t.Errorf("Expected LinkId %s, got %s", linkID, linkVisit.LinkId.String())
		}
		if linkVisit.VisitorId.String() != visitorId {
			t.Errorf("Expected VisitorId %s, got %s", visitorId, linkVisit.VisitorId.String())
		}
		if linkVisit.Ip.String() != ip {
			t.Errorf("Expected Ip %s, got %s", ip, linkVisit.Ip.String())
		}
		if string(linkVisit.UserAgent) != userAgent {
			t.Errorf("Expected UserAgent %s, got %s", userAgent, string(linkVisit.UserAgent))
		}
		if linkVisit.Id.String() == "" {
			t.Error("Expected ID to be generated")
		}
		if time.Time(linkVisit.CreatedOn).IsZero() {
			t.Error("Expected CreatedOn to be set")
		}
	})

	t.Run("Success with different IP formats", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name string
			ip   string
		}{
			{
				name: "IPv4",
				ip:   "192.168.1.1",
			},
			{
				name: "IPv4 with port",
				ip:   "192.168.1.1:8080",
			},
			{
				name: "IPv6",
				ip:   "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			},
			{
				name: "IPv6 short",
				ip:   "::1",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				linkID := "00000000-0000-0000-0000-000000000001"
				visitorId := "00000000-0000-0000-0000-000000000002"
				userAgent := "Mozilla/5.0"

				linkVisit, err := NewLinkVisit(linkID, visitorId, tt.ip, userAgent)
				if err != nil {
					t.Fatalf("Expected no error for %s, got %v", tt.name, err)
				}
				if linkVisit.Ip.String() == "" {
					t.Errorf("Expected Ip to be set for %s", tt.name)
				}
			})
		}
	})

	t.Run("Fails on invalid linkID", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name   string
			linkID string
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
				visitorId := "00000000-0000-0000-0000-000000000002"
				ip := "192.168.1.1"
				userAgent := "Mozilla/5.0"

				_, err := NewLinkVisit(tt.linkID, visitorId, ip, userAgent)
				if err == nil {
					t.Fatalf("Expected error for %s, got nil", tt.name)
				}
			})
		}
	})

	t.Run("Fails on invalid visitorId", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name      string
			visitorId string
		}{
			{
				name:      "invalid UUID format",
				visitorId: "invalid-uuid",
			},
			{
				name:      "empty string",
				visitorId: "",
			},
			{
				name:      "malformed UUID",
				visitorId: "00000000-0000-0000-0000",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				linkID := "00000000-0000-0000-0000-000000000001"
				ip := "192.168.1.1"
				userAgent := "Mozilla/5.0"

				_, err := NewLinkVisit(linkID, tt.visitorId, ip, userAgent)
				if err == nil {
					t.Fatalf("Expected error for %s, got nil", tt.name)
				}
			})
		}
	})

	t.Run("Fails on invalid IP", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name string
			ip   string
		}{
			{
				name: "empty string",
				ip:   "",
			},
			{
				name: "invalid format",
				ip:   "not-an-ip",
			},
			{
				name: "invalid IPv4",
				ip:   "256.256.256.256",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				linkID := "00000000-0000-0000-0000-000000000001"
				visitorId := "00000000-0000-0000-0000-000000000002"
				userAgent := "Mozilla/5.0"

				_, err := NewLinkVisit(linkID, visitorId, tt.ip, userAgent)
				if err == nil {
					t.Fatalf("Expected error for %s, got nil", tt.name)
				}
			})
		}
	})
}
