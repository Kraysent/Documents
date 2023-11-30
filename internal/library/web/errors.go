package web

import (
	"net/http"
)

const (
	DatabaseErrorCode      = "DATABASE_ERROR"
	InternalErrorCode      = "INTERNAL_ERROR"
	ValidationErrorCode    = "VALIDATION_ERROR"
	AuthorizationErrorCode = "AUTHORIZATION_ERROR"
)

type CodedError struct {
	Code     string `json:"code"`
	HTTPCode int    `json:"http_code"`
	Message  string `json:"message"`
}

func (e CodedError) Error() string {
	return e.Message
}

func ValidationError(err error) CodedError {
	return CodedError{
		Code:     ValidationErrorCode,
		HTTPCode: http.StatusBadRequest,
		Message:  err.Error(),
	}
}

func DatabaseError(err error) CodedError {
	return CodedError{
		Code:     DatabaseErrorCode,
		HTTPCode: http.StatusInternalServerError,
		Message:  err.Error(),
	}
}

func InternalError(err error) CodedError {
	return CodedError{
		Code:     InternalErrorCode,
		HTTPCode: http.StatusInternalServerError,
		Message:  err.Error(),
	}
}

func AuthorizationError(err error) CodedError {
	return CodedError{
		Code:     AuthorizationErrorCode,
		HTTPCode: http.StatusUnauthorized,
		Message:  err.Error(),
	}
}
