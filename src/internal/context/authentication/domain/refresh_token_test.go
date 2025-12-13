package domain

import (
	"testing"
	"time"
)

func TestNewRefreshToken(t *testing.T) {
	t.Run("Success on valid parameters", func(t *testing.T) {
		t.Parallel()
		userId := "00000000-0000-0000-0000-000000000001"
		token := "refresh-token-123"
		expiresAt := time.Now().Add(24 * time.Hour)

		refreshToken, err := NewRefreshToken(userId, token, expiresAt)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if refreshToken.UserId.String() != userId {
			t.Errorf("Expected UserId %s, got %s", userId, refreshToken.UserId.String())
		}
		if string(refreshToken.Token) != token {
			t.Errorf("Expected Token %s, got %s", token, string(refreshToken.Token))
		}
		if refreshToken.Id.String() == "" {
			t.Error("Expected ID to be generated")
		}
		if time.Time(refreshToken.CreatedOn).IsZero() {
			t.Error("Expected CreatedOn to be set")
		}
	})

	t.Run("Fails on invalid userId", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name   string
			userId string
		}{
			{
				name:   "invalid UUID format",
				userId: "invalid-uuid",
			},
			{
				name:   "empty string",
				userId: "",
			},
			{
				name:   "malformed UUID",
				userId: "00000000-0000-0000-0000",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				token := "refresh-token-123"
				expiresAt := time.Now().Add(24 * time.Hour)

				_, err := NewRefreshToken(tt.userId, token, expiresAt)
				if err == nil {
					t.Fatalf("Expected error for %s, got nil", tt.name)
				}
			})
		}
	})

	t.Run("Fails on empty token", func(t *testing.T) {
		t.Parallel()
		userId := "00000000-0000-0000-0000-000000000001"
		expiresAt := time.Now().Add(24 * time.Hour)

		_, err := NewRefreshToken(userId, "", expiresAt)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("Handles past expiration time", func(t *testing.T) {
		t.Parallel()
		userId := "00000000-0000-0000-0000-000000000001"
		token := "refresh-token-123"
		expiresAt := time.Now().Add(-24 * time.Hour)

		refreshToken, err := NewRefreshToken(userId, token, expiresAt)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if refreshToken == nil {
			t.Fatal("Expected refresh token to be created")
		}
	})
}

func TestRefreshTokenFromPrimitives(t *testing.T) {
	t.Run("Success on valid parameters", func(t *testing.T) {
		t.Parallel()
		id := "00000000-0000-0000-0000-000000000001"
		userId := "00000000-0000-0000-0000-000000000002"
		token := "refresh-token-123"
		expiresAt := time.Now().Add(24 * time.Hour)
		createdOn := time.Now().Add(-1 * time.Hour)

		refreshToken, err := RefreshTokenFromPrimitives(id, userId, token, expiresAt, createdOn)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if refreshToken.Id.String() != id {
			t.Errorf("Expected Id %s, got %s", id, refreshToken.Id.String())
		}
		if refreshToken.UserId.String() != userId {
			t.Errorf("Expected UserId %s, got %s", userId, refreshToken.UserId.String())
		}
		if string(refreshToken.Token) != token {
			t.Errorf("Expected Token %s, got %s", token, string(refreshToken.Token))
		}
	})

	t.Run("Fails on invalid id", func(t *testing.T) {
		t.Parallel()
		userId := "00000000-0000-0000-0000-000000000002"
		token := "refresh-token-123"
		expiresAt := time.Now().Add(24 * time.Hour)
		createdOn := time.Now()

		_, err := RefreshTokenFromPrimitives("invalid-uuid", userId, token, expiresAt, createdOn)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("Fails on invalid userId", func(t *testing.T) {
		t.Parallel()
		id := "00000000-0000-0000-0000-000000000001"
		token := "refresh-token-123"
		expiresAt := time.Now().Add(24 * time.Hour)
		createdOn := time.Now()

		_, err := RefreshTokenFromPrimitives(id, "invalid-uuid", token, expiresAt, createdOn)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("Fails on empty token", func(t *testing.T) {
		t.Parallel()
		id := "00000000-0000-0000-0000-000000000001"
		userId := "00000000-0000-0000-0000-000000000002"
		expiresAt := time.Now().Add(24 * time.Hour)
		createdOn := time.Now()

		_, err := RefreshTokenFromPrimitives(id, userId, "", expiresAt, createdOn)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}
