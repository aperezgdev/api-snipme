package middleware

import (
	"net/http"
	"strconv"
	"time"

	http_infra "github.com/aperezgdev/api-snipme/src/internal/context/shared/infrastructure/http"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latencies in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method", "status"},
	)
)

type PrometheusMiddleware struct{}

func NewPrometheusMiddleware() *PrometheusMiddleware {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	return &PrometheusMiddleware{}
}

func (PrometheusMiddleware) Handle(next http_infra.Route) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{w, http.StatusOK}
		next.Handler(rw, r)

		duration := time.Since(start).Seconds()

		path := r.URL.Path
		method := r.Method
		status := rw.statusCode

		httpRequestsTotal.WithLabelValues(path, method, strconv.Itoa(status)).Inc()
		httpRequestDuration.WithLabelValues(path, method, strconv.Itoa(status)).Observe(duration)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
