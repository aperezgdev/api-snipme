package domain

import shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"

type ShortLink struct {
	shared_domain.AggregateRoot
	Id            shared_domain.Id
	OriginalRoute ShortLinkOriginalRoute
	Code          ShortLinkCode
	Client        shared_domain.Id
	CreatedOn     shared_domain.CreatedOn
}

func NewShortLink(originalLink, client string) (*ShortLink, error) {
	idVO, err := shared_domain.NewID()
	originalPathVO, err := NewShortLinkOriginalRoute(originalLink)
	codeVO, err := NewCode()
	var clientVO shared_domain.Id
	if client != "" {
		clientVO, err = shared_domain.ParseID(client)
	}

	if err != nil {
		return nil, err
	}

	shortLink := &ShortLink{
		Id:            idVO,
		OriginalRoute: originalPathVO,
		Code:          codeVO,
		Client:        clientVO,
		CreatedOn:     shared_domain.NewCreatedOn(),
	}

	shortLink.Record(NewShortLinkCreated(idVO.String()))

	return shortLink, nil
}
