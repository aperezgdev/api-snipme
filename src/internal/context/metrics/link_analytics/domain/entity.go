package domain

import (
	link_country_view_counter_domain "github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_country_view_counter/domain"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/metrics/shared/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type LinkAnalytics struct {
	Id                        shared_domain_context.Id
	LinkId                    shared_domain_context.Id
	TotalViews                shared_domain.LinkViewsCounter
	UniqueViews               shared_domain.LinkViewsCounter
	LinkCountriesViewCounters []link_country_view_counter_domain.LinkCountryViewCounter
	UpdateOn                  shared_domain_context.UpdatedOn
}

func NewLinkAnalytics(linkID string) (*LinkAnalytics, error) {
	idVO, err := shared_domain_context.NewID()
	if err != nil {
		return nil, err
	}
	linkIDVO, err := shared_domain_context.ParseID(linkID)
	if err != nil {
		return nil, err
	}

	linkAnalytics := &LinkAnalytics{
		Id:                        idVO,
		LinkId:                    linkIDVO,
		TotalViews:                shared_domain.NewLinkViewsCounter(0),
		UniqueViews:               shared_domain.NewLinkViewsCounter(0),
		LinkCountriesViewCounters: []link_country_view_counter_domain.LinkCountryViewCounter{},
	}
	return linkAnalytics, nil
}
