package domain

import (
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type LinkVisit struct {
	shared_domain_context.AggregateRoot
	Id        shared_domain_context.Id
	LinkId    shared_domain_context.Id
	VisitorId shared_domain_context.Id
	Ip        LinkVisitIP
	UserAgent LinkVisitUserAgent
	CreatedOn shared_domain_context.CreatedOn
}

func NewLinkVisit(linkID, visitorId, ip, userAgent string) (*LinkVisit, error) {
	idVO, err := shared_domain_context.NewID()
	if err != nil {
		return nil, err
	}

	linkIDVO, err := shared_domain_context.ParseID(linkID)
	if err != nil {
		return nil, err
	}

	ipVO, err := NewLinkVisitIP(ip)
	if err != nil {
		return nil, err
	}

	visitorIdVO, err := shared_domain_context.ParseID(visitorId)
	if err != nil {
		return nil, err
	}

	linkVisit := &LinkVisit{
		Id:        idVO,
		LinkId:    linkIDVO,
		VisitorId: visitorIdVO,
		Ip:        ipVO,
		UserAgent: NewLinkVisitUserAgent(userAgent),
		CreatedOn: shared_domain_context.NewCreatedOn(),
	}

	return linkVisit, nil
}
