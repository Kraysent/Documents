package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"documents/internal/actions"
	"documents/internal/core"
	"documents/internal/entities"
)

func InsertDocument(request *http.Request, repo *core.Repository) (any, error) {
	ctx := context.Background()
	var document entities.Document

	data, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, entities.ValidationError(err)
	}

	if err := json.Unmarshal(data, &document); err != nil {
		return nil, entities.ValidationError(err)
	}

	if err := document.Validate(); err != nil {
		return nil, entities.ValidationError(err)
	}

	id, cErr := actions.InsertDocument(ctx, repo, document)
	if cErr != nil {
		return nil, cErr
	}

	return map[string]any{"id": id}, nil
}
