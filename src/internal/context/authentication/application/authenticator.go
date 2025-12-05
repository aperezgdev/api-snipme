package application

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/aperezgdev/api-snipme/src/internal/context/authentication/domain"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type AuthenticationResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
	User         *domain.User
}

type Authenticator struct {
	logger                shared_domain.Logger
	userRepo              domain.UserRepository
	refreshTokenRepo      domain.RefreshTokenRepository
	tokenManager          domain.TokenManager
	eventBus              shared_domain.EventBus
	userEmailUpdater      domain.UserEmailUpdater
	refreshTokenTTLDays   int
	jwtExpirantionMinutes int
}

func NewAuthenticator(
	logger shared_domain.Logger,
	userRepo domain.UserRepository,
	refreshTokenRepo domain.RefreshTokenRepository,
	tokenManager domain.TokenManager,
	eventBus shared_domain.EventBus,
	refreshTokenTTLDays int,
	jwtExpirantionMinutes int,
) *Authenticator {
	return &Authenticator{
		logger:                logger,
		userRepo:              userRepo,
		refreshTokenRepo:      refreshTokenRepo,
		tokenManager:          tokenManager,
		eventBus:              eventBus,
		userEmailUpdater:      *domain.NewUserEmailUpdater(logger, userRepo),
		refreshTokenTTLDays:   refreshTokenTTLDays,
		jwtExpirantionMinutes: jwtExpirantionMinutes,
	}
}

func (a *Authenticator) Run(
	ctx context.Context,
	provider domain.OAuthProvider,
	oauthSubject string,
	email string,
) (*AuthenticationResult, error) {
	a.logger.Info(ctx, "Authenticator - Run -  Authenticating user with OAuth", shared_domain.Field{Key: "provider", Value: provider.String()}, shared_domain.Field{Key: "email", Value: email})

	userOpt, err := a.userRepo.FindByOAuthProviderAndSubject(ctx, provider, oauthSubject)
	if err != nil {
		a.logger.Error(ctx, "Authenticator - Run - Failed to find user by OAuth provider and subject", shared_domain.Field{Key: "error", Value: err.Error()})
		return nil, err
	}

	var user *domain.User

	if userOpt.IsPresent() {
		a.logger.Info(ctx, "Authenticator - Run - User found, proceeding to authenticate", shared_domain.Field{Key: "user_id", Value: userOpt.Get().Id.String()})
		user = userOpt.Get()

		if string(user.Email) != email {
			a.logger.Info(ctx, "Authenticator - Run - User email has changed, updating email", shared_domain.Field{Key: "user_id", Value: user.Id.String()}, shared_domain.Field{Key: "old_email", Value: user.Email}, shared_domain.Field{Key: "new_email", Value: email})

			if err := a.userEmailUpdater.Run(ctx, user.Id.String(), email); err != nil {
				a.logger.Error(ctx, "Authenticator - Run - Failed to update user email", shared_domain.Field{Key: "error", Value: err.Error()})
				return nil, err
			}
		}
	} else {
		a.logger.Info(ctx, "Authenticator - Run - User not found, creating new user", shared_domain.Field{Key: "email", Value: email})

		user, err = domain.NewUser(email, provider, oauthSubject)
		if err != nil {
			return nil, err
		}

		if err := a.userRepo.Save(ctx, user); err != nil {
			a.logger.Error(ctx, "Authenticator - Run - Failed to save new user", shared_domain.Field{Key: "error", Value: err.Error()})
			return nil, err
		}

		for _, event := range user.PullDomainEvents() {
			a.eventBus.Publish(ctx, event)
		}
	}

	accessToken, err := a.tokenManager.Generate(user.Id.String(), string(user.Email))
	if err != nil {
		a.logger.Error(ctx, "Authenticator - Run - Failed to generate JWT", shared_domain.Field{Key: "error", Value: err.Error()})
		return nil, err
	}
	a.logger.Info(ctx, "Authenticator - Run - JWT generated successfully", shared_domain.Field{Key: "user_id", Value: user.Id.String()})

	refreshTokenStr, err := a.generateSecureToken()
	if err != nil {
		return nil, err
	}
	a.logger.Info(ctx, "Authenticator - Run - Refresh token generated successfully", shared_domain.Field{Key: "user_id", Value: user.Id.String()})

	refreshToken, err := domain.NewRefreshToken(
		user.Id.String(),
		refreshTokenStr,
		time.Now().Add(time.Duration(a.refreshTokenTTLDays)*24*time.Hour),
	)

	if err != nil {
		return nil, err
	}
	a.logger.Info(ctx, "Authenticator - Run - Refresh token created successfully", shared_domain.Field{Key: "user_id", Value: user.Id.String()})

	if err := a.refreshTokenRepo.Save(ctx, refreshToken); err != nil {
		a.logger.Error(ctx, "Authenticator - Run - Failed to save refresh token", shared_domain.Field{Key: "error", Value: err.Error()})
		return nil, err
	}
	a.logger.Info(ctx, "Authenticator - Run - Refresh token saved successfully", shared_domain.Field{Key: "user_id", Value: user.Id.String()})

	a.logger.Info(ctx, "Authenticator - Run - User authenticated successfully", shared_domain.Field{Key: "user_id", Value: user.Id.String()})

	return &AuthenticationResult{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		ExpiresIn:    a.jwtExpirantionMinutes * 60,
		User:         user,
	}, nil
}

func (a *Authenticator) generateSecureToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
