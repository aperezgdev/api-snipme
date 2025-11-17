package http

import (
	"errors"
	"net/http"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/application"
)

type DeleteShortLink struct {
	logger  shared_domain_context.Logger
	deleter application.ShortLinkRemover
}

func NewDeleteShortLinkHTTPHandler(
	logger shared_domain_context.Logger,
	deleter application.ShortLinkRemover,
) *DeleteShortLink {
	return &DeleteShortLink{
		logger:  logger,
		deleter: deleter,
	}
}

func (h *DeleteShortLink) Handler(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("short_code")

	err := h.deleter.Run(r.Context(), shortCode)
	var notFoundErr shared_domain_context.NotFoundError
	if err != nil && errors.As(err, &notFoundErr) {
		h.logger.Info(r.Context(), "DeleteShortLink - Handler - Short link not found", shared_domain_context.NewField("short_code", shortCode))
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if err != nil {
		h.logger.Error(r.Context(), "DeleteShortLink - Handler - Error trying to delete short link", shared_domain_context.NewField("error", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *DeleteShortLink) Method() string {
	return http.MethodDelete
}

func (h *DeleteShortLink) Route() string {
	return "/short-link/{short_code}"
}
