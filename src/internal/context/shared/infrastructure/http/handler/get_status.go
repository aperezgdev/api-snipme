package handler

import "net/http"

type GetStatusHTTPHandler struct{}

func NewGetStatusHTTPHandler() *GetStatusHTTPHandler {
	return &GetStatusHTTPHandler{}
}

func (h *GetStatusHTTPHandler) Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
	w.WriteHeader(http.StatusOK)
}

func (h *GetStatusHTTPHandler) Method() string {
	return http.MethodGet
}

func (h *GetStatusHTTPHandler) Route() string {
	return "/status"
}
