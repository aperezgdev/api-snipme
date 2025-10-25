package middleware

import (
	"net/http"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	http_infra "github.com/aperezgdev/api-snipme/src/internal/context/shared/infrastructure/http"
)

type RecoveryMiddleware struct {
	logger shared_domain_context.Logger
}

func NewRecoveryMiddleware(logger shared_domain_context.Logger) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		logger: logger,
	}
}

func (rm RecoveryMiddleware) Handle(next http_infra.Route) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				rm.logger.Error(r.Context(), "Recovered from panic", shared_domain_context.NewField("error", err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.Handler(w, r)
	})
}
