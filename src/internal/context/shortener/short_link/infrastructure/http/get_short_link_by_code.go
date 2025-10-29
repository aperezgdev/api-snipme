package http

import (
	"errors"
	"net/http"

	application_link_visit "github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_visit/application"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/application"
)

type GetShortLinkByCodeHTTPHandler struct {
	logger  shared_domain_context.Logger
	finder  application.ShortLinkFinderByCode
	creator application_link_visit.LinkVisitCreator
}

func NewGetShortLinkByCodeHTTPHandler(
	logger shared_domain_context.Logger,
	finder application.ShortLinkFinderByCode,
	creator application_link_visit.LinkVisitCreator,
) *GetShortLinkByCodeHTTPHandler {
	return &GetShortLinkByCodeHTTPHandler{
		logger:  logger,
		finder:  finder,
		creator: creator,
	}
}

func (h *GetShortLinkByCodeHTTPHandler) Handler(w http.ResponseWriter, req *http.Request) {
	code := req.PathValue("code")

	shortLink, err := h.finder.Run(req.Context(), code)
	var validationErr shared_domain_context.ValidationError
	if errors.As(err, &validationErr) {
		h.logger.Error(req.Context(), "ProjectHttpHandler - Error creating project", shared_domain_context.NewField("error", err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var notFoundErr shared_domain_context.NotFoundError
	if errors.As(err, &notFoundErr) {
		h.logger.Error(req.Context(), "GetShortLinkByCodeHTTPHandler - Short link not found", shared_domain_context.NewField("code", code))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		h.logger.Error(req.Context(), "PostShortLinkHTTPHanlder - Internal error", shared_domain_context.NewField("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = h.creator.Run(req.Context(), shortLink.Id.String(), req.RemoteAddr, req.UserAgent())
	if errors.Is(err, shared_domain_context.ValidationError{}) {
		h.logger.Error(req.Context(), "GetShortLinkByCodeHTTPHandler - Error creating link visit", shared_domain_context.NewField("error", err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if err != nil {
		h.logger.Error(req.Context(), "PostShortLinkHTTPHanlder - Internal error", shared_domain_context.NewField("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, string(shortLink.OriginalRoute), http.StatusFound)
}

func (h *GetShortLinkByCodeHTTPHandler) Route() string {
	return "/{code}"
}

func (h *GetShortLinkByCodeHTTPHandler) Method() string {
	return http.MethodGet
}
