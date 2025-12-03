package domain

type TokenClaims struct {
	UserID string
	Email  string
}

type TokenManager interface {
	Generate(userID, email string) (string, error)
	Validate(tokenString string) (*TokenClaims, error)
}
