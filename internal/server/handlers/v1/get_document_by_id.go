package v1

import (
	"context"
	"net/http"

	"documents/internal/actions"
	"documents/internal/actions/schema"
	"documents/internal/core"
	"documents/internal/library/web"
	gschema "github.com/gorilla/schema"
)

func GetDocumentByID(r *http.Request, repo *core.Repository) (any, error) {
	ctx := context.Background()
	var request schema.GetDocumentByIDRequest
	decoder := gschema.NewDecoder()

	if err := decoder.Decode(&request, r.URL.Query()); err != nil {
		return nil, web.ValidationError(err)
	}

	if err := request.Validate(); err != nil {
		return nil, web.ValidationError(err)
	}

	document, cErr := actions.GetDocumentByID(ctx, repo, request)
	if cErr != nil {
		return nil, cErr
	}

	return document, nil
}
