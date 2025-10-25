package middleware

import (
	"context"
	"net/http"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	http_infra "github.com/aperezgdev/api-snipme/src/internal/context/shared/infrastructure/http"
	"github.com/google/uuid"
)

type RequestIDMiddleware struct {
	logger shared_domain_context.Logger
}

func NewRequestIDMiddleware(logger shared_domain_context.Logger) *RequestIDMiddleware {
	return &RequestIDMiddleware{
		logger: logger,
	}
}

func (rm RequestIDMiddleware) Handle(next http_infra.Route) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		w.Header().Set("X-Request-ID", requestID)
		rm.logger.Info(r.Context(), "Request ID set", shared_domain_context.NewField("request_id", requestID), shared_domain_context.NewField("method", r.Method), shared_domain_context.NewField("url", r.URL.String()))
		r = r.WithContext(context.WithValue(r.Context(), "request_id", requestID))
		next.Handler(w, r)
	})
}
