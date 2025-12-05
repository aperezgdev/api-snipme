package domain

import "net/mail"

type Email string

func NewEmail(email string) (Email, error) {
	err := validateEmail(email)
	if err != nil {
		return "", err
	}
	return Email(email), nil
}

func validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return NewValidationError("email", "must be a valid email address")
	}
	return err
}
