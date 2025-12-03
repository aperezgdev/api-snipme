package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

type OAuthUserInfo struct {
	Email   string
	Subject string
}

type OAuthClient interface {
	GetAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error)
	GetUserInfo(ctx context.Context, token *oauth2.Token) (*OAuthUserInfo, error)
}

type GoogleOAuthClient struct {
	config *oauth2.Config
}

func NewGoogleOAuthClient(clientID, clientSecret, redirectURL string) *GoogleOAuthClient {
	return &GoogleOAuthClient{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func (c *GoogleOAuthClient) GetAuthURL(state string) string {
	return c.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (c *GoogleOAuthClient) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return c.config.Exchange(ctx, code)
}

func (c *GoogleOAuthClient) GetUserInfo(ctx context.Context, token *oauth2.Token) (*OAuthUserInfo, error) {
	client := c.config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("failed to get user info from Google")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo struct {
		Email string `json:"email"`
		ID    string `json:"id"`
	}

	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &OAuthUserInfo{
		Email:   userInfo.Email,
		Subject: userInfo.ID,
	}, nil
}

type GitHubOAuthClient struct {
	config *oauth2.Config
}

func NewGitHubOAuthClient(clientID, clientSecret, redirectURL string) *GitHubOAuthClient {
	return &GitHubOAuthClient{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
		},
	}
}

func (c *GitHubOAuthClient) GetAuthURL(state string) string {
	return c.config.AuthCodeURL(state)
}

func (c *GitHubOAuthClient) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return c.config.Exchange(ctx, code)
}

func (c *GitHubOAuthClient) GetUserInfo(ctx context.Context, token *oauth2.Token) (*OAuthUserInfo, error) {
	client := c.config.Client(ctx, token)
	
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get user info from GitHub: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo struct {
		ID    int64  `json:"id"`
		Email string `json:"email"`
	}

	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	if userInfo.Email == "" {
		emailResp, err := client.Get("https://api.github.com/user/emails")
		if err != nil {
			return nil, err
		}
		defer emailResp.Body.Close()

		emailBody, err := io.ReadAll(emailResp.Body)
		if err != nil {
			return nil, err
		}

		var emails []struct {
			Email    string `json:"email"`
			Primary  bool   `json:"primary"`
			Verified bool   `json:"verified"`
		}

		if err := json.Unmarshal(emailBody, &emails); err != nil {
			return nil, err
		}

		for _, e := range emails {
			if e.Primary && e.Verified {
				userInfo.Email = e.Email
				break
			}
		}

		if userInfo.Email == "" && len(emails) > 0 {
			userInfo.Email = emails[0].Email
		}
	}

	return &OAuthUserInfo{
		Email:   userInfo.Email,
		Subject: fmt.Sprintf("%d", userInfo.ID),
	}, nil
}

type OAuthClientMock struct {
	mock.Mock
}

func (m *OAuthClientMock) GetAuthURL(state string) string {
	args := m.Called(state)
	return args.String(0)
}

func (m *OAuthClientMock) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*oauth2.Token), args.Error(1)
}

func (m *OAuthClientMock) GetUserInfo(ctx context.Context, token *oauth2.Token) (*OAuthUserInfo, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*OAuthUserInfo), args.Error(1)
}
