package domain

import (
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

const LinkViewCounterIncrementedEventName = "LinkViewCounterUpdated"

type LinkViewCounterIncremented struct {
	shared_domain_context.DomainEventBase
	TotalViews  uint
	UniqueViews uint
}

func NewLinkViewCounterIncrementedDomainEvent(linkId string, countryCode string, totalViews, uniqueViews uint) LinkViewCounterIncremented {
	return LinkViewCounterIncremented{
		DomainEventBase: shared_domain_context.NewDomainEvent(
			linkId,
			LinkViewCounterIncrementedEventName,
		),
		TotalViews:  totalViews,
		UniqueViews: uniqueViews,
	}
}

func (LinkViewCounterIncremented) Name() string {
	return LinkViewCounterIncrementedEventName
}
