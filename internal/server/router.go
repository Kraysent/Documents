package server

import (
	"net/http"

	"documents/internal/actions"
	"documents/internal/actions/schema"
	"documents/internal/core"
	"documents/internal/server/handlers/admin"
	"documents/internal/server/handlers/auth"
	v1 "documents/internal/server/handlers/v1"
)

func GetAuthHandlers() []AuthHandler {
	return []AuthHandler{
		{
			Path:     "/api/auth/google/login",
			Method:   http.MethodGet,
			Function: auth.GoogleLoginHandler,
		},
		{
			Path:     "/api/auth/google/callback",
			Method:   http.MethodGet,
			Function: auth.GoogleCallbackHandler,
		},
	}
}

func GetHandlers() []CommonHandler {
	return []CommonHandler{
		{
			Path:   "/api/ping",
			Method: http.MethodGet,
			Function: func(request *http.Request, repository *core.Repository) (any, error) {
				return map[string]any{"ping": "pong"}, nil
			},
		},
		{ // TODO: remove after MVP is done
			Path:     "/api/internal/token",
			Method:   http.MethodPost,
			Function: admin.IssueToken,
		},
		{
			Path:   "/api/v1/document",
			Method: http.MethodPost,
			Function: v1.JSONHandler[
				schema.InsertDocumentRequest, schema.InsertDocumentResponse,
			](actions.InsertDocument),
		},
		{
			Path:   "/api/v1/document/id",
			Method: http.MethodGet,
			Function: v1.QueryHandler[
				schema.GetDocumentByIDRequest, schema.GetDocumentResponse,
			](actions.GetDocumentByID),
		},
		{
			Path:   "/api/v1/document",
			Method: http.MethodDelete,
			Function: v1.QueryHandler[
				schema.DeleteDocumentRequest, schema.DeleteDocumentResponse,
			](actions.DeleteDocument),
		},
		{
			Path:   "/api/v1/user/documents",
			Method: http.MethodGet,
			Function: v1.QueryHandler[
				schema.GetUserDocumentsRequest, schema.GetUserDocumentsResponse,
			](actions.GetUserDocuments),
		},
		{
			Path:   "/api/v1/link",
			Method: http.MethodPost,
			Function: v1.JSONHandler[
				schema.CreateLinkRequest, schema.CreateLinkResponse,
			](actions.CreateLink),
		},
		{
			Path:   "/api/v1/link",
			Method: http.MethodGet,
			Function: v1.QueryHandler[
				schema.GetDocumentByLinkRequest, schema.GetDocumentResponse,
			](actions.GetDocumentByLink),
		},
		{
			Path:   "/api/v1/link",
			Method: http.MethodDelete,
			Function: v1.QueryHandler[
				schema.DisableLinkRequest, schema.DisableLinkResponse,
			](actions.DisableLink),
		},
		{
			Path:   "/api/v1/links",
			Method: http.MethodGet,
			Function: v1.QueryHandler[
				schema.GetLinksRequest, schema.GetLinksResponse,
			](actions.GetLinks),
		},
	}
}
