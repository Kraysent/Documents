package validation

import (
	"fmt"
)

type ValidationError struct {
	messages []string
}

func NewValidationError(value any, condition string) ValidationError {
	return ValidationError{
		messages: []string{
			fmt.Sprintf("field with value '%v' does not satisfy the condition '%s'", value, condition),
		},
	}
}

func mergeErrors(errs []ValidationError) ValidationError {
	resError := ValidationError{}

	for _, err := range errs {
		resError.messages = append(resError.messages, err.messages...)
	}

	return resError
}

func (e ValidationError) Error() string {
	msg := "validations errors: "

	for _, message := range e.messages {
		msg += fmt.Sprintf("%s; ", message)
	}

	return msg
}

func (e ValidationError) IsEmpty() bool {
	return len(e.messages) == 0
}
