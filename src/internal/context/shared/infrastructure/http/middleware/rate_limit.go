package middleware

import (
	"net/http"
	"sync"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	http_infra "github.com/aperezgdev/api-snipme/src/internal/context/shared/infrastructure/http"
	"golang.org/x/time/rate"
)

type RateLimitMiddleware struct {
	logger shared_domain_context.Logger
	ips    map[string]*rate.Limiter
	mutex  sync.Mutex
	rate   rate.Limit
	burst  int
}

func NewRateLimitMiddleware(logger shared_domain_context.Logger, r rate.Limit, burst int) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		logger: logger,
		ips:    make(map[string]*rate.Limiter),
		rate:   r,
		burst:  burst,
	}
}

func (rlm *RateLimitMiddleware) AddIP(ip string) *rate.Limiter {
	rlm.mutex.Lock()
	defer rlm.mutex.Unlock()

	limiter := rate.NewLimiter(rlm.rate, int(rlm.burst))
	rlm.ips[ip] = limiter
	return limiter
}

func (rlm *RateLimitMiddleware) GetLimiter(ip string) *rate.Limiter {
	rlm.mutex.Lock()
	limiter, exists := rlm.ips[ip]
	rlm.mutex.Unlock()

	if !exists {
		return rlm.AddIP(ip)
	}
	return limiter
}

func (rlm *RateLimitMiddleware) Handle(next http_infra.Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		limiter := rlm.GetLimiter(ip)

		if !limiter.Allow() {
			rlm.logger.Info(r.Context(), "Rate limit exceeded", shared_domain_context.NewField("ip", ip))
			w.WriteHeader(429)
			w.Write([]byte("Too Many Requests"))
			return
		}

		next.Handler(w, r)
	}
}
