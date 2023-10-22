package v1

import (
	"encoding/json"
	"io"
	"net/http"

	"documents/internal/actions"
	"documents/internal/actions/schema"
	"documents/internal/core"
	"documents/internal/library/web"
)

func InsertDocument(r *http.Request, repo *core.Repository) (any, error) {
	var request schema.InsertDocumentRequest

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, web.ValidationError(err)
	}

	if err := json.Unmarshal(data, &request); err != nil {
		return nil, web.ValidationError(err)
	}

	if err := request.Validate(); err != nil {
		return nil, web.ValidationError(err)
	}

	response, cErr := actions.InsertDocument(r.Context(), repo, request)
	if cErr != nil {
		return nil, cErr
	}

	return response, nil
}
