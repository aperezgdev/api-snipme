package http

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aperezgdev/api-snipme/src/internal/context/authentication/application"
	"github.com/aperezgdev/api-snipme/src/internal/context/authentication/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/authentication/infrastructure"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type GetOAuthCallbackHandler struct {
	logger        shared_domain.Logger
	oauthClient   infrastructure.OAuthClient
	authenticator *application.Authenticator
	provider      domain.OAuthProvider
	stateSecret   string
}

func NewGetOAuthCallbackHandler(
	logger shared_domain.Logger,
	oauthClient infrastructure.OAuthClient,
	authenticator *application.Authenticator,
	provider domain.OAuthProvider,
	stateSecret string,
) *GetOAuthCallbackHandler {
	return &GetOAuthCallbackHandler{
		logger:        logger,
		oauthClient:   oauthClient,
		authenticator: authenticator,
		provider:      provider,
		stateSecret:   stateSecret,
	}
}

func (h *GetOAuthCallbackHandler) Handler(w http.ResponseWriter, req *http.Request) {
	h.logger.Info(req.Context(), "GetOAuthCallbackHandler - Handler - Handling OAuth callback", shared_domain.NewField("provider", h.provider.String()))
	code := req.URL.Query().Get("code")
	if code == "" {
		h.logger.Error(req.Context(), "GetOAuthCallbackHandler - Handler - Missing code parameter")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "missing code parameter"})
		return
	}

	receivedState := req.URL.Query().Get("state")
	if receivedState == "" {
		h.logger.Error(req.Context(), "GetOAuthCallbackHandler - Handler - Missing state parameter")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "missing state parameter"})
		return
	}

	stateCookie, err := req.Cookie("oauth_state")
	if err != nil {
		h.logger.Error(req.Context(), "GetOAuthCallbackHandler - Handler - Missing state cookie")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "missing state cookie - possible CSRF attack"})
		return
	}

	parts := strings.Split(stateCookie.Value, ".")
	if len(parts) != 2 {
		h.logger.Error(req.Context(), "GetOAuthCallbackHandler - Handler - Invalid state cookie format")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid state format"})
		return
	}
	expectedState, expectedSignature := parts[0], parts[1]

	mac := hmac.New(sha256.New, []byte(h.stateSecret))
	mac.Write([]byte(expectedState))
	validSignature := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expectedSignature), []byte(validSignature)) {
		h.logger.Error(req.Context(), "GetOAuthCallbackHandler - Handler - Invalid state signature - possible CSRF attack")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid state signature"})
		return
	}

	if receivedState != expectedState {
		h.logger.Error(req.Context(), "GetOAuthCallbackHandler - Handler - State mismatch - possible CSRF attack")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "state mismatch - CSRF detected"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "oauth_state",
		Value:  "",
		Path:   "/auth",
		MaxAge: -1,
	})

	h.logger.Info(req.Context(), "GetOAuthCallbackHandler - Handler - Exchanging code for token")
	token, err := h.oauthClient.ExchangeCode(req.Context(), code)
	if err != nil {
		h.logger.Error(req.Context(), "GetOAuthCallbackHandler - Handler - Failed to exchange code for token", shared_domain.NewField("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to exchange code"})
		return
	}
	h.logger.Info(req.Context(), "GetOAuthCallbackHandler - Handler - Retrieving user info from OAuth provider")

	userInfo, err := h.oauthClient.GetUserInfo(req.Context(), token)
	if err != nil {
		h.logger.Error(req.Context(), "GetOAuthCallbackHandler - Handler - Failed to get user info", shared_domain.NewField("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to get user info"})
		return
	}
	h.logger.Info(req.Context(), "GetOAuthCallbackHandler - Handler - Authenticating user", shared_domain.NewField("user_subject", userInfo.Subject))

	result, err := h.authenticator.Run(
		req.Context(),
		h.provider,
		userInfo.Subject,
		userInfo.Email,
	)
	if err != nil {
		h.logger.Error(req.Context(), "GetOAuthCallbackHandler - Handler - Authentication failed", shared_domain.NewField("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "authentication failed"})
		return
	}
	h.logger.Info(req.Context(), "GetOAuthCallbackHandler - Handler - Successfully authenticated user", shared_domain.NewField("user_id", result.User.Id.String()))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"token_type":    "Bearer",
		"expires_in":    result.ExpiresIn,
		"user": map[string]string{
			"id":    result.User.Id.String(),
			"email": string(result.User.Email),
		},
	})
}

func (h *GetOAuthCallbackHandler) Route() string {
	return "/auth/" + h.provider.String() + "/callback"
}

func (h *GetOAuthCallbackHandler) Method() string {
	return http.MethodGet
}
