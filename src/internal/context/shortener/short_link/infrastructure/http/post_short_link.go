package http

import (
	"encoding/json"
	"errors"
	"net/http"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/application"
)

type PostShortLinkHTTPHanlder struct {
	logger  shared_domain_context.Logger
	creator application.ShortLinkCreator
}

type postShortLinkHttpRequest struct {
	OriginalURL string `json:"original_url"`
	ClientID    string `json:"client_id"`
}

func NewPostShortLinkHTTPHandler(logger shared_domain_context.Logger, creator application.ShortLinkCreator) *PostShortLinkHTTPHanlder {
	return &PostShortLinkHTTPHanlder{
		logger:  logger,
		creator: creator,
	}
}

func (h *PostShortLinkHTTPHanlder) Handler(w http.ResponseWriter, req *http.Request) {
	request := postShortLinkHttpRequest{}
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		h.logger.Error(req.Context(), "PostShortLinkHTTPHanlder - Invalid JSON", shared_domain_context.NewField("error", err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request body"))
		return
	}

	_, err := h.creator.Run(req.Context(), request.OriginalURL, request.ClientID)

	var validationErr shared_domain_context.ValidationError
	if err != nil && errors.As(err, &validationErr) {
		h.logger.Error(req.Context(), "PostShortLinkHTTPHanlder - Error creating short link", shared_domain_context.NewField("error", err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if err != nil {
		h.logger.Error(req.Context(), "PostShortLinkHTTPHanlder - Internal error", shared_domain_context.NewField("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *PostShortLinkHTTPHanlder) Route() string {
	return "/short-links"
}

func (h *PostShortLinkHTTPHanlder) Method() string {
	return http.MethodPost
}
