package web

import (
	"net/http"
)

const (
	DatabaseErrorCode      = "database_error"
	InternalErrorCode      = "internal_error"
	ValidationErrorCode    = "validation_error"
	AuthorizationErrorCode = "authorization_error"
	NotFoundErrorCode      = "not_found"
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

func NotFoundError(err error) CodedError {
	return CodedError{
		Code:     NotFoundErrorCode,
		HTTPCode: http.StatusNotFound,
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
