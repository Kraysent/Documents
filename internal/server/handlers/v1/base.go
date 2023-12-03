package v1

import (
	"encoding/json"
	"io"
	"net/http"

	"documents/internal/actions"
	"documents/internal/core"
	"documents/internal/library/web"
	gschema "github.com/gorilla/schema"
)

type Handler func(r *http.Request, repo *core.Repository) (any, error)

func QueryHandler[RequestType actions.Request, ResponseType any](
	action actions.Action[RequestType, ResponseType],
) Handler {
	return func(r *http.Request, repo *core.Repository) (any, error) {
		var request RequestType
		decoder := gschema.NewDecoder()

		if err := decoder.Decode(&request, r.URL.Query()); err != nil {
			return nil, web.ValidationError(err)
		}

		if err := request.Validate(); err != nil {
			return nil, web.ValidationError(err)
		}

		document, cErr := action(r.Context(), repo, request)
		if cErr != nil {
			return nil, cErr
		}

		return document, nil
	}
}

func JSONHandler[RequestType actions.Request, ResponseType any](
	action actions.Action[RequestType, ResponseType],
) Handler {
	return func(r *http.Request, repo *core.Repository) (any, error) {
		var request RequestType

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

		response, err := action(r.Context(), repo, request)
		if err != nil {
			return nil, err
		}

		return response, nil
	}
}
