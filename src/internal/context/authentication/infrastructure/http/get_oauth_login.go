package http

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/http"

	"github.com/aperezgdev/api-snipme/src/internal/context/authentication/infrastructure"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type GetOAuthLoginHandler struct {
	logger      shared_domain.Logger
	oauthClient infrastructure.OAuthClient
	provider    string
	stateSecret string
}

func NewGetOAuthLoginHandler(
	logger shared_domain.Logger,
	oauthClient infrastructure.OAuthClient,
	provider string,
	stateSecret string,
) *GetOAuthLoginHandler {
	return &GetOAuthLoginHandler{
		logger:      logger,
		oauthClient: oauthClient,
		provider:    provider,
		stateSecret: stateSecret,
	}
}

func (h *GetOAuthLoginHandler) Handler(w http.ResponseWriter, req *http.Request) {
	h.logger.Info(req.Context(), "GetOAuthLoginHandler - Handler - Generating OAuth login URL", shared_domain.NewField("provider", h.provider))
	stateBytes := make([]byte, 32)
	if _, err := rand.Read(stateBytes); err != nil {
		h.logger.Error(req.Context(), "GetOAuthLoginHandler - Handler - Failed to generate state parameter", shared_domain.NewField("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	state := base64.URLEncoding.EncodeToString(stateBytes)
	h.logger.Info(req.Context(), "GetOAuthLoginHandler - Handler - Generated state parameter", shared_domain.NewField("state", state))

	mac := hmac.New(sha256.New, []byte(h.stateSecret))
	mac.Write([]byte(state))
	signature := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state + "." + signature,
		Path:     "/auth",
		MaxAge:   600,
		HttpOnly: true,
		Secure:   req.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})
	h.logger.Info(req.Context(), "GetOAuthLoginHandler - Handler - Set cookie")

	authURL := h.oauthClient.GetAuthURL(state)

	h.logger.Info(req.Context(), "GetOAuthLoginHandler - Handler - Redirecting to OAuth provider")

	http.Redirect(w, req, authURL, http.StatusTemporaryRedirect)
}

func (h *GetOAuthLoginHandler) Route() string {
	return "/auth/" + h.provider + "/login"
}

func (h *GetOAuthLoginHandler) Method() string {
	return http.MethodGet
}
