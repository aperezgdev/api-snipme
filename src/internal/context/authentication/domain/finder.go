package domain

import (
	"context"

	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type UserFinder struct {
	logger   shared_domain.Logger
	userRepo UserRepository
}

func NewUserFinder(
	logger shared_domain.Logger,
	userRepo UserRepository,
) *UserFinder {
	return &UserFinder{
		logger:   logger,
		userRepo: userRepo,
	}
}

func (f *UserFinder) Run(ctx context.Context, userId string) (*User, error) {
	f.logger.Info(ctx, "UserFinder - Run - Finding user by ID", shared_domain.Field{Key: "user_id", Value: userId})

	parsedId, err := shared_domain.ParseID(userId)
	if err != nil {
		f.logger.Error(ctx, "UserFinder - Run - Invalid user ID format", shared_domain.Field{Key: "error", Value: err.Error()})
		return nil, shared_domain.NewValidationError("user_id", "invalid format")
	}

	userOpt, err := f.userRepo.FindById(ctx, parsedId)
	if err != nil {
		f.logger.Error(ctx, "UserFinder - Run - Failed to find user by ID", shared_domain.Field{Key: "error", Value: err.Error()})
		return nil, err
	}

	if !userOpt.IsPresent() {
		f.logger.Info(ctx, "UserFinder - Run - User not found", shared_domain.Field{Key: "user_id", Value: userId})
		return nil, shared_domain.NewNotFoundError("user not found")
	}

	f.logger.Info(ctx, "UserFinder - Run - User found", shared_domain.Field{Key: "user_id", Value: userId})
	return userOpt.Get(), nil
}