package domain

import domain_shared "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"

type Client struct {
	Id        domain_shared.Id
	Name      ClientName
	Email     ClientEmail
	CreatedOn domain_shared.CreatedOn
}

func NewClient(name, email string) (*Client, error) {
	idVO, err := domain_shared.NewID()
	if err != nil {
		return &Client{}, err
	}
	nameVO, err := NewClientName(name)
	if err != nil {
		return &Client{}, err
	}
	emailVO, err := NewClientEmail(email)
	if err != nil {
		return &Client{}, err
	}
	return &Client{
		Id:        idVO,
		Name:      nameVO,
		Email:     emailVO,
		CreatedOn: domain_shared.NewCreatedOn(),
	}, nil
}
