package models

import "fmt"

type NotEnoughArgsError struct {
	error
	Message string
}

func NewNotEnoughArgsError(message string) *NotEnoughArgsError {
	return &NotEnoughArgsError{Message: message}
}

func (r *NotEnoughArgsError) Error() string {
	return fmt.Sprintf(fmt.Sprintf("Not enough arguments: %s", r.Message))
}
