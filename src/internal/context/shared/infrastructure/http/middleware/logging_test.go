package middleware

import (
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	http "github.com/aperezgdev/api-snipme/src/internal/context/shared/infrastructure/http"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestLoggerMiddleware(t *testing.T) {
	logger := &shared_domain_context.DummyLogger{}
	mw := NewLoggerMiddleware(logger)
	rec := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	mw.Handle(http.DummyRoute{}).ServeHTTP(rec, r)
	assert.Equal(t, 200, rec.Code)
}
