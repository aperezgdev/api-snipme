package middleware

import (
	"net/http/httptest"
	"testing"

	customhttp "github.com/aperezgdev/api-snipme/src/internal/context/shared/infrastructure/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestPrometheusMiddleware(t *testing.T) {
	prometheus.Unregister(httpRequestsTotal)
	prometheus.Unregister(httpRequestDuration)
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

	mw := NewPrometheusMiddleware()
	rec := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)

	mw.Handle(customhttp.DummyRoute{}).ServeHTTP(rec, r)

	assert.Equal(t, 200, rec.Code)

	count := testutil.ToFloat64(httpRequestsTotal.WithLabelValues("/test", "GET", "200"))
	assert.Equal(t, 1.0, count)
}
