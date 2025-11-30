package domain

import (
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

const LinkVisitsProcessedEventName = "LinkVisitsProcessed"

type LinkVisitsProcessed struct {
    shared_domain_context.DomainEventBase
    TotalViews  uint
    UniqueViews uint
}

func NewLinkVisitsProcessedDomainEvent(linkId string, totalViews, uniqueViews uint) LinkVisitsProcessed {
    return LinkVisitsProcessed{
        DomainEventBase: shared_domain_context.NewDomainEvent(
            linkId,
            LinkVisitsProcessedEventName,
        ),
        TotalViews:  totalViews,
        UniqueViews: uniqueViews,
    }
}

func (LinkVisitsProcessed) Name() string {
    return LinkVisitsProcessedEventName
}