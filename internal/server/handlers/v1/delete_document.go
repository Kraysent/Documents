package v1

import (
	"net/http"

	"documents/internal/actions"
	"documents/internal/actions/schema"
	"documents/internal/core"
	"documents/internal/library/web"
	gschema "github.com/gorilla/schema"
)

func DeleteDocument(r *http.Request, repo *core.Repository) (any, error) {
	var request schema.DeleteDocumentRequest
	decoder := gschema.NewDecoder()

	if err := decoder.Decode(&request, r.URL.Query()); err != nil {
		return nil, web.ValidationError(err)
	}

	if err := request.Validate(); err != nil {
		return nil, web.ValidationError(err)
	}

	document, cErr := actions.DeleteDocument(r.Context(), repo, request)
	if cErr != nil {
		return nil, cErr
	}

	return document, nil
}
