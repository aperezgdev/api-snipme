package domain

import domain_shared "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"

type ClientName string

func NewClientName(name string) (ClientName, error) {
	err := validateName(name)
	if err != nil {
		return "", err
	}
	return ClientName(name), nil
}

func validateName(name string) error {
	if len(name) == 0 {
		return domain_shared.NewValidationError("name", "Name cannot be empty")
	}
	if len(name) > 50 {
		return domain_shared.NewValidationError("name", "Name cannot be longer than 50 characters")
	}
	return nil
}
