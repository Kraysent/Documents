package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"documents/internal/actions"
	"documents/internal/actions/schema"
	"documents/internal/core"
)

func InsertDocument(r *http.Request, repo *core.Repository) (any, error) {
	ctx := context.Background()
	var request schema.InsertDocumentRequest

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, actions.ValidationError(err)
	}

	if err := json.Unmarshal(data, &request); err != nil {
		return nil, actions.ValidationError(err)
	}

	if err := request.Validate(); err != nil {
		return nil, actions.ValidationError(err)
	}

	response, cErr := actions.InsertDocument(ctx, repo, request)
	if cErr != nil {
		return nil, cErr
	}

	return response, nil
}
