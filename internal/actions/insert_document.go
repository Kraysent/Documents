package actions

import (
	"context"
	"fmt"

	"documents/internal/actions/schema"
	"documents/internal/core"
	"documents/internal/library/web"
)

func InsertDocument(
	ctx context.Context, repo *core.Repository, r schema.InsertDocumentRequest,
) (*schema.InsertDocumentResponse, error) {
	userID := repo.SessionManager.GetInt64(ctx, "user_id")
	if userID == 0 {
		return nil, web.AuthorizationError(fmt.Errorf("failed to authorize"))
	}

	id, err := repo.Storages.Documents.AddDocument(ctx, userID, r.Type, r.Attributes)
	if err != nil {
		return nil, web.DatabaseError(err)
	}

	return &schema.InsertDocumentResponse{ID: id}, nil
}
