package application

import (
	"context"

	refresh_token_domain "github.com/aperezgdev/api-snipme/src/internal/context/authentication/refresh_token/domain"
	user_domain "github.com/aperezgdev/api-snipme/src/internal/context/authentication/user/domain"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type TokenValidator struct {
	logger       shared_domain.Logger
	tokenManager refresh_token_domain.TokenManager
	userRepo     user_domain.UserRepository
}

func NewTokenValidator(
	logger shared_domain.Logger,
	tokenManager refresh_token_domain.TokenManager,
	userRepo user_domain.UserRepository,
) *TokenValidator {
	return &TokenValidator{
		logger:       logger,
		tokenManager: tokenManager,
		userRepo:     userRepo,
	}
}

func (v *TokenValidator) Validate(ctx context.Context, token string) (*user_domain.User, error) {
	claims, err := v.tokenManager.Validate(token)
	if err != nil {
		v.logger.Debug(ctx, "Invalid JWT token", shared_domain.Field{Key: "error", Value: err.Error()})
		return nil, shared_domain.NewValidationError("token", "invalid or expired")
	}

	userId, err := shared_domain.ParseID(claims.UserID)
	if err != nil {
		return nil, err
	}

	userOpt, err := v.userRepo.FindById(ctx, userId)
	if err != nil {
		v.logger.Error(ctx, "Failed to find user by ID", shared_domain.Field{Key: "error", Value: err.Error()})
		return nil, err
	}

	if !userOpt.IsPresent() {
		return nil, shared_domain.NewNotFoundError("user not found")
	}

	return userOpt.Get(), nil
}
