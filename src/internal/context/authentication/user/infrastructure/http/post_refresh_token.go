package http

import (
	"encoding/json"
	"errors"
	"net/http"

	refresh_token_application "github.com/aperezgdev/api-snipme/src/internal/context/authentication/refresh_token/application"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type PostRefreshTokenHandler struct {
	logger         shared_domain.Logger
	tokenRefresher *refresh_token_application.TokenRefresher
}

func NewPostRefreshTokenHandler(
	logger shared_domain.Logger,
	tokenRefresher *refresh_token_application.TokenRefresher,
) *PostRefreshTokenHandler {
	return &PostRefreshTokenHandler{
		logger:         logger,
		tokenRefresher: tokenRefresher,
	}
}

type requestBody struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *PostRefreshTokenHandler) Handler(w http.ResponseWriter, req *http.Request) {
	h.logger.Info(req.Context(), "PostRefreshTokenHandler - Handler - Handling token refresh request")
	var requestBody requestBody
	if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
		h.logger.Error(req.Context(), "PostRefreshTokenHandler - Handler - Failed to decode request body", shared_domain.NewField("error", err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	if requestBody.RefreshToken == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error(req.Context(), "PostRefreshTokenHandler - Handler - Missing refresh_token in request body")
		json.NewEncoder(w).Encode(map[string]string{"error": "refresh_token is required"})
		return
	}

	result, err := h.tokenRefresher.Refresh(req.Context(), requestBody.RefreshToken)
	if err != nil {
		var validationErr shared_domain.ValidationError
		var notFoundErr shared_domain.NotFoundError

		if errors.As(err, &validationErr) || errors.As(err, &notFoundErr) {
			w.WriteHeader(http.StatusUnauthorized)
			h.logger.Error(req.Context(), "PostRefreshTokenHandler - Handler - Invalid or expired refresh token", shared_domain.NewField("error", err))
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid or expired refresh token"})
			return
		}

		h.logger.Error(req.Context(), "PostRefreshTokenHandler - Handler - Failed to refresh token", shared_domain.NewField("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to refresh token"})
		return
	}
	h.logger.Info(req.Context(), "PostRefreshTokenHandler - Handler - Successfully refreshed access token")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token": result.AccessToken,
		"token_type":   "Bearer",
		"expires_in":   result.ExpiresIn,
	})
}

func (h *PostRefreshTokenHandler) Route() string {
	return "/auth/refresh"
}

func (h *PostRefreshTokenHandler) Method() string {
	return http.MethodPost
}
