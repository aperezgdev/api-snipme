package application

import (
	"context"

	"github.com/aperezgdev/api-snipme/src/internal/context/authentication/domain"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type RefreshTokenResult struct {
	AccessToken string
	ExpiresIn   int
}

type TokenRefresher struct {
	logger           shared_domain.Logger
	refreshTokenRepo domain.RefreshTokenRepository
	userRepo         domain.UserRepository
	tokenManager     domain.TokenManager
	jwtTokenExpirationMinutes int
}

func NewTokenRefresher(
	logger shared_domain.Logger,
	refreshTokenRepo domain.RefreshTokenRepository,
	userRepo domain.UserRepository,
	tokenManager domain.TokenManager,
	jwtTokenExpirationMinutes int,
) *TokenRefresher {
	return &TokenRefresher{
		logger:           logger,
		refreshTokenRepo: refreshTokenRepo,
		userRepo:         userRepo,
		tokenManager:     tokenManager,
		jwtTokenExpirationMinutes: jwtTokenExpirationMinutes,
	}
}

func (r *TokenRefresher) Refresh(ctx context.Context, refreshTokenStr string) (*RefreshTokenResult, error) {
	r.logger.Info(ctx, "Refreshing access token")

	tokenOpt, err := r.refreshTokenRepo.FindByToken(ctx, refreshTokenStr)
	if err != nil {
		r.logger.Error(ctx, "Failed to find refresh token", shared_domain.Field{Key: "error", Value: err.Error()})
		return nil, err
	}

	if !tokenOpt.IsPresent() {
		return nil, shared_domain.NewNotFoundError("refresh token not found")
	}

	token := tokenOpt.Get()

	if token.IsExpired() {
		r.logger.Info(ctx, "Refresh token is expired")
		return nil, shared_domain.NewValidationError("refresh_token", "expired")
	}

	userOpt, err := r.userRepo.FindById(ctx, token.UserId)
	if err != nil {
		r.logger.Error(ctx, "Failed to find user", shared_domain.Field{Key: "error", Value: err.Error()})
		return nil, err
	}

	if !userOpt.IsPresent() {
		return nil, shared_domain.NewNotFoundError("user not found")
	}

	user := userOpt.Get()

	accessToken, err := r.tokenManager.Generate(user.Id.String(), string(user.Email))
	if err != nil {
		r.logger.Error(ctx, "Failed to generate JWT", shared_domain.Field{Key: "error", Value: err.Error()})
		return nil, err
	}

	r.logger.Info(ctx, "Access token refreshed successfully", shared_domain.Field{Key: "user_id", Value: user.Id.String()})

	return &RefreshTokenResult{
		AccessToken: accessToken,
		ExpiresIn:   r.jwtTokenExpirationMinutes * 60,
	}, nil
}
