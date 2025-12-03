package application

import (
	"context"

	auth_domain "github.com/aperezgdev/api-snipme/src/internal/context/authentication/domain"
	domain_shared "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/client/domain"
)

type CreatorOnUserCreated struct {
	logger domain_shared.Logger
	repo   domain.ClientRepository
}

func NewCreatorOnUserCreated(logger domain_shared.Logger, repo domain.ClientRepository) *CreatorOnUserCreated {
	return &CreatorOnUserCreated{
		logger: logger,
		repo:   repo,
	}
}

func (c CreatorOnUserCreated) On(ctx context.Context, event domain_shared.DomainEvent) error {
	userCreatedEvent, ok := event.(auth_domain.UserCreated)
	if !ok {
		c.logger.Error(ctx, "CreatorOnUserCreated - On - Invalid event type", 
			domain_shared.NewField("expected", auth_domain.UserCreatedEventName), 
			domain_shared.NewField("received", event.Name()))
		return nil
	}

	c.logger.Info(ctx, "CreatorOnUserCreated - On - Creating client for new user",
		domain_shared.NewField("user_id", userCreatedEvent.AggregateID()),
		domain_shared.NewField("email", userCreatedEvent.Email()),
		domain_shared.NewField("name", userCreatedEvent.UserName()))

	client, err := domain.NewClient(userCreatedEvent.UserName(), userCreatedEvent.Email())
	if err != nil {
		c.logger.Error(ctx, "CreatorOnUserCreated - On - Error creating client domain entity",
			domain_shared.NewField("error", err.Error()),
			domain_shared.NewField("user_id", userCreatedEvent.AggregateID()))
		return err
	}

	err = c.repo.Save(ctx, *client)
	if err != nil {
		c.logger.Error(ctx, "CreatorOnUserCreated - On - Error saving client",
			domain_shared.NewField("error", err.Error()),
			domain_shared.NewField("user_id", userCreatedEvent.AggregateID()))
		return err
	}

	c.logger.Info(ctx, "CreatorOnUserCreated - On - Client created successfully for user",
		domain_shared.NewField("user_id", userCreatedEvent.AggregateID()),
		domain_shared.NewField("client_id", client.Id.String()))

	return nil
}

func (c CreatorOnUserCreated) SubscribedTo() []string {
	return []string{auth_domain.UserCreatedEventName}
}
