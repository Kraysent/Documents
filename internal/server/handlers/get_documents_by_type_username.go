package handlers

import (
	"context"
	"net/http"

	"documents/internal/actions"
	"documents/internal/actions/schema"
	"documents/internal/core"
	gschema "github.com/gorilla/schema"
)

func GetDocumentsByUsernameAndType(r *http.Request, repo *core.Repository) (any, error) {
	ctx := context.Background()
	var request schema.GetDocumentByUsernameAndTypeRequest
	decoder := gschema.NewDecoder()

	if err := decoder.Decode(&request, r.URL.Query()); err != nil {
		return nil, actions.ValidationError(err)
	}

	if err := request.Validate(); err != nil {
		return nil, actions.ValidationError(err)
	}

	document, cErr := actions.GetDocumentsByUsernameAndType(ctx, repo, request)
	if cErr != nil {
		return nil, cErr
	}

	return document, nil
}
