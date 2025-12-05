package infrastructure

import (
	"errors"
	"time"

	"github.com/aperezgdev/api-snipme/src/internal/context/authentication/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret            []byte
	expirationMinutes int
}

func NewJWTManager(secret string, expirationMinutes int) *JWTManager {
	return &JWTManager{
		secret:            []byte(secret),
		expirationMinutes: expirationMinutes,
	}
}

func (m *JWTManager) Generate(userID, email string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(m.expirationMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *JWTManager) Validate(tokenString string) (*domain.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return &domain.TokenClaims{
			UserID: claims.UserID,
			Email:  claims.Email,
		}, nil
	}

	return nil, errors.New("invalid token")
}

type JWTManagerMock struct {
	mock.Mock
}

func (m *JWTManagerMock) Generate(userID, email string) (string, error) {
	args := m.Called(userID, email)
	return args.String(0), args.Error(1)
}

func (m *JWTManagerMock) Validate(tokenString string) (*domain.TokenClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TokenClaims), args.Error(1)
}
