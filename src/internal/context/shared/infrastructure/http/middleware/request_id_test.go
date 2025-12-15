package middleware

import (
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	http "github.com/aperezgdev/api-snipme/src/internal/context/shared/infrastructure/http"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestRequestIDMiddleware(t *testing.T) {
	logger := &shared_domain_context.DummyLogger{}
	mw := NewRequestIDMiddleware(logger)
	rec := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	mw.Handle(http.DummyRoute{}).ServeHTTP(rec, r)
	assert.NotEmpty(t, rec.Header().Get("X-Request-ID"))
}
