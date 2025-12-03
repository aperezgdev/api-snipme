package test

import (
	"database/sql"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func generateTestToken(userId string) (string, error) {
	jwtSecret := "test-jwt-secret-key-for-testing"
	expirationMinutes := 60

	claims := JWTClaims{
		UserID: userId,
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirationMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func createTestUser(db *sql.DB) (string, error) {
	userId := uuid.New()
	email := "test@example.com"
	provider := "google"
	subject := "test-oauth-subject-" + time.Now().Format("20060102150405")

	query := `
		INSERT INTO "user" (id, email, oauth_provider, oauth_subject, created_on)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (email) DO UPDATE SET oauth_subject = EXCLUDED.oauth_subject
		RETURNING id
	`

	var returnedId string
	err := db.QueryRow(query, userId.String(), email, provider, subject).Scan(&returnedId)
	if err != nil {
		return "", err
	}

	return returnedId, nil
}
