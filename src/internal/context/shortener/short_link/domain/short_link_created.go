package domain

import shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"

const ShortLinkCreatedEventName = "ShortLinkCreated"

type ShortLinkCreated struct {
	shared_domain_context.DomainEventBase
}

func NewShortLinkCreated(shortLinkID string) ShortLinkCreated {
	return ShortLinkCreated{
		DomainEventBase: shared_domain_context.NewDomainEvent(
			shortLinkID,
			ShortLinkCreatedEventName,
		),
	}
}

func (e ShortLinkCreated) Name() string {
	return ShortLinkCreatedEventName
}
