package http

import (
	"encoding/json"
	"errors"
	"net/http"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/application"
)

type postPublicShortLinkRequest struct {
	OriginalLink string `json:"original_link"`
}

type PostPublicShortLinkHTTPHandler struct {
	logger shared_domain_context.Logger
	creator application.PublicShortLinkCreator
}

func NewPostPublicShortLinkHTTPHandler(
	logger shared_domain_context.Logger,
	creator application.PublicShortLinkCreator,
) *PostPublicShortLinkHTTPHandler {
	return &PostPublicShortLinkHTTPHandler{
		logger: logger,
		creator: creator,
	}
}

func (h *PostPublicShortLinkHTTPHandler) Handler(w http.ResponseWriter, req *http.Request) {
	h.logger.Info(req.Context(), "PostPublicShortLinkHTTPHandler - Handler - Start")
	var request postPublicShortLinkRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		h.logger.Error(req.Context(), "PostPublicShortLinkHTTPHandler - Handler - Invalid JSON", shared_domain_context.NewField("error", err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request body"))
		return
	}
	h.logger.Info(req.Context(), "PostPublicShortLinkHTTPHandler - Handler - Request parsed", shared_domain_context.NewField("original_link", request.OriginalLink))

	if request.OriginalLink == "" {
		h.logger.Error(req.Context(), "PostPublicShortLinkHTTPHandler - Handler - Original link is required")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("original_link is required"))
		return
	}

	_, err := h.creator.Run(req.Context(), request.OriginalLink)
	if err != nil && errors.As(err, &shared_domain_context.ValidationError{}) {
		h.logger.Error(req.Context(), "PostPublicShortLinkHTTPHandler - Handler - Validation error", shared_domain_context.NewField("error", err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	
	if err != nil {
		h.logger.Error(req.Context(), "PostPublicShortLinkHTTPHandler - Handler - Error creating public short link", shared_domain_context.NewField("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.logger.Info(req.Context(), "PostPublicShortLinkHTTPHandler - Handler - Public short link created successfully")

	w.WriteHeader(http.StatusCreated)
}

func (*PostPublicShortLinkHTTPHandler) Method() string {
	return http.MethodPost
}

func (*PostPublicShortLinkHTTPHandler) Route() string {
	return "/public-short-links"
}