package domain

import (
	"context"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type UserEmailUpdater struct {
	logger     shared_domain_context.Logger
	userRepo   UserRepository
	finderById UserFinder
}

func NewUserEmailUpdater(
	logger shared_domain_context.Logger,
	userRepo UserRepository,
) *UserEmailUpdater {
	return &UserEmailUpdater{
		logger:     logger,
		userRepo:   userRepo,
		finderById: *NewUserFinder(logger, userRepo),
	}
}

func (u *UserEmailUpdater) Run(ctx context.Context, userId, newEmail string) error {
	u.logger.Info(ctx, "UserEmailUpdater - Run - Starting email update", shared_domain_context.NewField("user_id", userId), shared_domain_context.NewField("new_email", newEmail))

	user, err := u.finderById.Run(ctx, userId)
	if err != nil {
		u.logger.Error(ctx, "UserEmailUpdater - Run - Failed to find user by ID", shared_domain_context.NewField("user_id", userId), shared_domain_context.NewField("error", err.Error()))
		return err
	}

	emailVO, err := shared_domain_context.NewEmail(newEmail)
	if err != nil {
		u.logger.Error(ctx, "UserEmailUpdater - Run - Invalid email format", shared_domain_context.NewField("new_email", newEmail), shared_domain_context.NewField("error", err.Error()))
		return err
	}

	user.Email = emailVO

	if err := u.userRepo.Update(ctx, user); err != nil {
		u.logger.Error(ctx, "UserEmailUpdater - Run - Failed to update user email", shared_domain_context.NewField("user_id", userId), shared_domain_context.NewField("error", err.Error()))
		return err
	}

	u.logger.Info(ctx, "UserEmailUpdater - Run - Successfully updated user email", shared_domain_context.NewField("user_id", userId))
	return nil

}
