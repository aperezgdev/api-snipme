package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/aperezgdev/api-snipme/src/internal/context/authentication/application"
	"github.com/aperezgdev/api-snipme/src/internal/context/authentication/domain"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	shared_http "github.com/aperezgdev/api-snipme/src/internal/context/shared/infrastructure/http"
)

type contextKey string

const AuthenticatedUserKey contextKey = "authenticated_user"

type AuthenticationMiddleware struct {
	logger         shared_domain.Logger
	tokenValidator *application.TokenValidator
}

func NewAuthenticationMiddleware(
	logger shared_domain.Logger,
	tokenValidator *application.TokenValidator,
) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		logger:         logger,
		tokenValidator: tokenValidator,
	}
}

func (m *AuthenticationMiddleware) Handle(next shared_http.Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.logger.Info(r.Context(), "AuthenticationMiddleware - Handle - Authenticating request")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.logger.Debug(r.Context(), "AuthenticationMiddleware - Handle - Missing Authorization header")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing Authorization header"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.logger.Debug(r.Context(), "AuthenticationMiddleware - Handle - Invalid Authorization header format")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid Authorization header format"))
			return
		}

		token := parts[1]

		user, err := m.tokenValidator.Validate(r.Context(), token)
		if err != nil {
			m.logger.Debug(r.Context(), "AuthenticationMiddleware - Handle - Invalid token", shared_domain.NewField("err", err))
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid or expired token"))
			return
		}

		m.logger.Info(r.Context(), "AuthenticationMiddleware - Handle - Successfully authenticated user", shared_domain.NewField("user_id", user.Id.String()))
		ctx := context.WithValue(r.Context(), AuthenticatedUserKey, user)
		next.Handler(w, r.WithContext(ctx))
	}
}

func GetAuthenticatedUser(ctx context.Context) (*domain.User, bool) {
	user, ok := ctx.Value(AuthenticatedUserKey).(*domain.User)
	return user, ok
}
