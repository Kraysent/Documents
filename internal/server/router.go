package server

import (
	"net/http"

	"documents/internal/server/handlers"
)

func GetHandlers() []CommonHandler {
	return []CommonHandler{
		{
			Path:     "/api/v1/document",
			Method:   http.MethodPost,
			Function: handlers.InsertDocument,
		},
		{
			Path:     "/api/v1/document",
			Method:   http.MethodGet,
			Function: handlers.GetDocumentsByUsernameAndType,
		},
		{
			Path:     "/api/v1/document/id",
			Method:   http.MethodGet,
			Function: handlers.GetDocumentByID,
		},
		{
			Path:     "/api/v1/document",
			Method:   http.MethodDelete,
			Function: handlers.DeleteDocument,
		},
		{
			Path:     "/api/v1/user/documents",
			Method:   http.MethodGet,
			Function: handlers.GetUserDocuments,
		},
		{
			Path:     "/api/v1/user",
			Method:   http.MethodPost,
			Function: handlers.CreateUser,
		},
	}
}
