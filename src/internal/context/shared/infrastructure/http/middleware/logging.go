package middleware

import (
	"net/http"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	http_infra "github.com/aperezgdev/api-snipme/src/internal/context/shared/infrastructure/http"
)

type LoggerMiddleware struct {
	logger shared_domain_context.Logger
}

func NewLoggerMiddleware(logger shared_domain_context.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{
		logger: logger,
	}
}

func (lm LoggerMiddleware) Handle(next http_infra.Route) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lm.logger.Info(r.Context(), "Request received", shared_domain_context.NewField("method", r.Method), shared_domain_context.NewField("url", r.URL.String()), shared_domain_context.NewField("remote_addr", r.RemoteAddr))
		next.Handler(w, r)
		lm.logger.Info(r.Context(), "Request completed", shared_domain_context.NewField("method", r.Method), shared_domain_context.NewField("url", r.URL.String()), shared_domain_context.NewField("remote_addr", r.RemoteAddr))
	})
}
