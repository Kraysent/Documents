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

func CreateUser(request *http.Request, repo *core.Repository) (any, error) {
	ctx := context.Background()
	var r schema.CreateUserRequest

	data, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, actions.ValidationError(err)
	}

	if err := json.Unmarshal(data, &r); err != nil {
		return nil, actions.ValidationError(err)
	}

	if err := r.Validate(); err != nil {
		return nil, actions.ValidationError(err)
	}

	response, cErr := actions.CreateUser(ctx, repo, r)
	if cErr != nil {
		return nil, cErr
	}

	return response, nil
}
