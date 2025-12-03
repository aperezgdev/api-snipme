package application

import (
	"context"

	"github.com/aperezgdev/api-snipme/src/internal/context/authentication/domain"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type TokenValidator struct {
	logger     shared_domain.Logger
	tokenManager domain.TokenManager
	userRepo   domain.UserRepository
}

func NewTokenValidator(
	logger shared_domain.Logger,
	tokenManager domain.TokenManager,
	userRepo domain.UserRepository,
) *TokenValidator {
	return &TokenValidator{
		logger:     logger,
		tokenManager: tokenManager,
		userRepo:   userRepo,
	}
}

func (v *TokenValidator) Validate(ctx context.Context, token string) (*domain.User, error) {
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
