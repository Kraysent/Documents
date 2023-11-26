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

func CreateLink(r *http.Request, repo *core.Repository) (any, error) {
	var request schema.CreateLinkRequest

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

	response, cErr := actions.CreateLink(r.Context(), repo, request)
	if cErr != nil {
		return nil, cErr
	}

	return response, nil
}
