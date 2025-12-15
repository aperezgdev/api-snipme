package domain

import (
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

const LinkVisitsProcessedEventName = "LinkVisitsProcessed"

type IpsVisits struct {
	Ip          string
	TotalViews  uint
	UniqueViews uint
}

type LinkVisitsProcessed struct {
	shared_domain_context.DomainEventBase
	TotalViews   uint
	UniqueViews  uint
	IpsFrequency []IpsVisits
}

func NewLinkVisitsProcessedDomainEvent(linkId string, totalViews, uniqueViews uint, ipsFrequency []IpsVisits) LinkVisitsProcessed {
	return LinkVisitsProcessed{
		DomainEventBase: shared_domain_context.NewDomainEvent(
			linkId,
			LinkVisitsProcessedEventName,
		),
		TotalViews:   totalViews,
		UniqueViews:  uniqueViews,
		IpsFrequency: ipsFrequency,
	}
}

func (LinkVisitsProcessed) Name() string {
	return LinkVisitsProcessedEventName
}
