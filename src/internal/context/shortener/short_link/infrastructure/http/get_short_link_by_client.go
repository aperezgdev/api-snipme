package http

import (
	"errors"
	"net/http"

	"encoding/json"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/application"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
)

type GetShortLinkByClientHTTPHandler struct {
	logger shared_domain_context.Logger
	finder application.ShortLinkFinderByClient
}

func NewGetShortLinkByClientHTTPHandler(
	logger shared_domain_context.Logger,
	finder application.ShortLinkFinderByClient,
) *GetShortLinkByClientHTTPHandler {
	return &GetShortLinkByClientHTTPHandler{
		logger: logger,
		finder: finder,
	}
}

type shortLinkResponse struct {
	ID        string `json:"id"`
	ShortCode string `json:"short_code"`
	URL       string `json:"url"`
	ClientID  string `json:"client_id"`
}

func (h *GetShortLinkByClientHTTPHandler) Handler(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	h.logger.Info(r.Context(), "GetShortLinkByClient - Handler - Received request", shared_domain_context.NewField("client_id", clientID))

	shortLinks, err := h.finder.Run(r.Context(), clientID)
	var notFoundErr shared_domain_context.NotFoundError
	if errors.As(err, &notFoundErr) {
		h.logger.Info(r.Context(), "GetShortLinkByClient - Handler - No short links found for client", shared_domain_context.NewField("client_id", clientID))
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if err != nil {
		h.logger.Error(r.Context(), "GetShortLinkByClient - Handler - Error trying to find short links", shared_domain_context.NewField("error", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	shortLinkReponses := pkg.Map(shortLinks, func(sl *domain.ShortLink) shortLinkResponse {
		return shortLinkResponse{
			ID:        sl.Id.String(),
			ShortCode: string(sl.Code),
			URL:       string(sl.OriginalRoute),
			ClientID:  sl.Client.String(),
		}
	})

	json, err := json.Marshal(shortLinkReponses)
	if err != nil {
		h.logger.Error(r.Context(), "GetShortLinkByClient - Handler - Error trying to marshal response", shared_domain_context.NewField("error", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(json)
	if err != nil {
		h.logger.Error(r.Context(), "GetShortLinkByClient - Handler - Error trying to write response", shared_domain_context.NewField("error", err.Error()))
	}
}

func (h *GetShortLinkByClientHTTPHandler) Method() string {
	return http.MethodGet
}

func (h *GetShortLinkByClientHTTPHandler) Route() string {
	return "/short-link"
}
