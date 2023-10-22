package admin

import (
	"encoding/json"
	"io"
	"net/http"

	"documents/internal/core"
	"documents/internal/library/web"
)

type IssueTokenRequest struct {
	UserID int64 `json:"user_id"`
}

func IssueToken(r *http.Request, repo *core.Repository) (any, error) {
	var request IssueTokenRequest

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, web.ValidationError(err)
	}

	if err := json.Unmarshal(data, &request); err != nil {
		return nil, web.ValidationError(err)
	}

	repo.SessionManager.Put(r.Context(), "user_id", request.UserID)

	return nil, nil
}
