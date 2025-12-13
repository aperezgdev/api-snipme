package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type PanicRoute struct{}

func (PanicRoute) Handler(w http.ResponseWriter, r *http.Request) {
	panic("test panic")
}

func (PanicRoute) Route() string {
	return "/"
}

func (PanicRoute) Method() string {
	return "GET"
}

func TestRecoveryMiddleware(t *testing.T) {
	logger := &shared_domain_context.DummyLogger{}
	mw := NewRecoveryMiddleware(logger)
	rec := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	panicRoute := PanicRoute{}
	mw.Handle(panicRoute).ServeHTTP(rec, r)
	assert.Equal(t, 500, rec.Code)
}
