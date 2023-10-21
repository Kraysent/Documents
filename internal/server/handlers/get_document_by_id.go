package handlers

import (
	"context"
	"net/http"

	"documents/internal/actions"
	schema2 "documents/internal/actions/schema"
	"documents/internal/core"
	"documents/internal/entities"
	"github.com/gorilla/schema"
)

func GetDocumentByID(r *http.Request, repo *core.Repository) (any, error) {
	ctx := context.Background()
	var request schema2.GetDocumentByIDRequest
	decoder := schema.NewDecoder()

	if err := decoder.Decode(&request, r.URL.Query()); err != nil {
		return nil, entities.ValidationError(err)
	}

	if err := request.Validate(); err != nil {
		return nil, entities.ValidationError(err)
	}

	document, cErr := actions.GetDocumentByID(ctx, repo, request)
	if cErr != nil {
		return nil, cErr
	}

	return document, nil
}