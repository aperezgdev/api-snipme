package middleware

import (
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	http "github.com/aperezgdev/api-snipme/src/internal/context/shared/infrastructure/http"
)

func TestRateLimitMiddleware(t *testing.T) {
	logger := &shared_domain_context.DummyLogger{}
	mw := NewRateLimitMiddleware(logger, 1, 1)
	rec := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	h := mw.Handle(http.DummyRoute{})
	h.ServeHTTP(rec, r)
	assert.Equal(t, 200, rec.Code)
	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, r)
	assert.Equal(t, 429, rec2.Code)
}
