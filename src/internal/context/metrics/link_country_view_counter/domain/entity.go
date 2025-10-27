package domain

import (
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/metrics/shared/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type LinkCountryViewCounter struct {
	shared_domain_context.AggregateRoot
	Id          shared_domain_context.Id
	LinkId      shared_domain_context.Id
	CountryCode CountryCode
	TotalViews  shared_domain.LinkViewsCounter
	UniqueViews shared_domain.LinkViewsCounter
	CreatedOn   shared_domain_context.CreatedOn
}

func NewLinkCountryViewCounter(linkId, countryCode string) (*LinkCountryViewCounter, error) {
	idVO, err := shared_domain_context.NewID()
	if err != nil {
		return nil, err
	}
	linkID, err := shared_domain_context.ParseID(linkId)
	if err != nil {
		return nil, err
	}
	countryCodeVO, err := NewCountryCode(countryCode)
	if err != nil {
		return nil, err
	}

	linkCountryViewCounter := &LinkCountryViewCounter{
		Id:          idVO,
		LinkId:      linkID,
		CountryCode: countryCodeVO,
		TotalViews:  0,
		UniqueViews: 0,
	}

	return linkCountryViewCounter, nil
}
func (lcs *LinkCountryViewCounter) Increment(totalViews, uniqueViews uint) *LinkCountryViewCounter {
	lcs.TotalViews += shared_domain.LinkViewsCounter(totalViews)
	lcs.UniqueViews += shared_domain.LinkViewsCounter(uniqueViews)
	lcs.Record(NewLinkViewCounterIncrementedDomainEvent(
		lcs.LinkId.String(),
		string(lcs.CountryCode),
		uint(lcs.TotalViews),
		uint(lcs.UniqueViews),
	))
	return lcs
}
