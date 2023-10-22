package server

import (
	"net/http"

	"documents/internal/server/handlers/admin"
	v1 "documents/internal/server/handlers/v1"
)

func GetHandlers() []CommonHandler {
	return []CommonHandler{
		{
			Path:     "/api/internal/token",
			Method:   http.MethodPost,
			Function: admin.IssueToken,
		},
		{
			Path:     "/api/v1/document",
			Method:   http.MethodPost,
			Function: v1.InsertDocument,
		},
		{
			Path:     "/api/v1/document",
			Method:   http.MethodGet,
			Function: v1.GetDocumentsByUsernameAndType,
		},
		{
			Path:     "/api/v1/document/id",
			Method:   http.MethodGet,
			Function: v1.GetDocumentByID,
		},
		{
			Path:     "/api/v1/document",
			Method:   http.MethodDelete,
			Function: v1.DeleteDocument,
		},
		{
			Path:     "/api/v1/user/documents",
			Method:   http.MethodGet,
			Function: v1.GetUserDocuments,
		},
	}
}
