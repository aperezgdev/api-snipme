package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_analytics/application"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type GetLinkAnalyticsByLinkHTTPHandler struct {
	logger shared_domain_context.Logger
	finder application.LinkAnalyticsFinder
}

func NewGetLinkAnalyticsByLinkHTTPHandler(
	logger shared_domain_context.Logger,
	finder application.LinkAnalyticsFinder,
) *GetLinkAnalyticsByLinkHTTPHandler {
	return &GetLinkAnalyticsByLinkHTTPHandler{
		logger: logger,
		finder: finder,
	}
}

type linkAnalyticsResponse struct {
	LinkId      string `json:"linkId"`
	TotalViews  uint   `json:"totalViews"`
	UniqueViews uint   `json:"uniqueViews"`
}

func (h *GetLinkAnalyticsByLinkHTTPHandler) Handler(w http.ResponseWriter, r *http.Request) {
	linkID := r.URL.Query().Get("link_id")
	h.logger.Info(r.Context(), "GetLinkAnalyticsByLink - Handler - Received request", shared_domain_context.NewField("link_id", linkID))

	if linkID == "" {
		h.logger.Info(r.Context(), "GetLinkAnalyticsByLink - Handler - link_id is missing in query parameters")
		http.Error(w, "Bad Request: missing link_id", http.StatusBadRequest)
		return
	}

	linkAnalytics, err := h.finder.Run(r.Context(), linkID)
	var notFoundErr shared_domain_context.NotFoundError
	if err != nil && errors.As(err, &notFoundErr) {
		h.logger.Info(r.Context(), "GetLinkAnalyticsByLink - Handler - No link analytics found for link", shared_domain_context.NewField("link_id", linkID))
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	var validationErr shared_domain_context.ValidationError
	if err != nil && errors.As(err, &validationErr) {
		h.logger.Info(r.Context(), "GetLinkAnalyticsByLink - Handler - Validation error", shared_domain_context.NewField("error", err.Error()))
		http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		h.logger.Error(r.Context(), "GetLinkAnalyticsByLink - Handler - Error trying to find link analytics", shared_domain_context.NewField("error", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := linkAnalyticsResponse{
		LinkId:      linkAnalytics.LinkId.String(),
		TotalViews:  uint(linkAnalytics.TotalViews),
		UniqueViews: uint(linkAnalytics.UniqueViews),
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.logger.Error(r.Context(), "GetLinkAnalyticsByLink - Handler - Error encoding response", shared_domain_context.NewField("error", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.logger.Info(r.Context(), "GetLinkAnalyticsByLink - Handler - Successfully responded", shared_domain_context.NewField("link_id", linkID))
}

func (GetLinkAnalyticsByLinkHTTPHandler) Method() string {
	return http.MethodGet
}

func (GetLinkAnalyticsByLinkHTTPHandler) Route() string {
	return "/link-analytics"
}
