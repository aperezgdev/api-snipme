package domain

import (
	"net/mail"

	domain_shared "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type ClientEmail string

func NewClientEmail(email string) (ClientEmail, error) {
	err := validateEmail(email)
	if err != nil {
		return "", err
	}
	return ClientEmail(email), nil
}

func validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return domain_shared.NewValidationError("email", "must be a valid email address")
	}
	return err
}
