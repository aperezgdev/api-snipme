package domain

import "fmt"

type ValidationError struct {
	Field   string
	Message string
}

type NotFoundError struct {
	Message string
}

func NewValidationError(field, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Message: message,
	}
}

func NewNotFoundError(message string) NotFoundError {
	return NotFoundError{
		Message: message,
	}
}

func (err NotFoundError) Error() string {
	return fmt.Sprintf("Not Found: %s", err.Message)
}

func (err ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", err.Field, err.Message)
}
