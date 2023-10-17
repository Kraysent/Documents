package handlers

import (
	"context"
	"net/http"

	"documents/internal/actions"
	"documents/internal/core"
	"documents/internal/entities"
	schema2 "documents/internal/server/schema"

	"github.com/gorilla/schema"
)

func GetDocumentByUsernameAndType(r *http.Request, repo *core.Repository) (any, error) {
	ctx := context.Background()
	var request schema2.GetDocumentByUsernameAndTypeRequest
	decoder := schema.NewDecoder()

	if err := decoder.Decode(&request, r.URL.Query()); err != nil {
		return nil, entities.ValidationError(err)
	}

	if err := request.Validate(); err != nil {
		return nil, entities.ValidationError(err)
	}

	document, cErr := actions.GetDocumentByUsernameAndType(ctx, repo, request.Username, request.Type)
	if cErr != nil {
		return nil, cErr
	}

	return document, nil
}
