package actions

import (
	"context"
	"fmt"

	"documents/internal/actions/schema"
	"documents/internal/core"
	"documents/internal/library/web"
	"documents/internal/storage/documents"
)

func InsertDocument(
	ctx context.Context, repo *core.Repository, r schema.InsertDocumentRequest,
) (*schema.InsertDocumentResponse, error) {
	userID := repo.SessionManager.GetInt64(ctx, "user_id")
	if userID == 0 {
		return nil, web.AuthorizationError(fmt.Errorf("failed to authorize"))
	}

	result, err := repo.Storages.Documents.AddDocument(ctx, documents.AddDocumentRequest{
		Name:        r.Name,
		Owner:       userID,
		Description: r.Description,
	})
	if err != nil {
		return nil, web.DatabaseError(err)
	}

	return &schema.InsertDocumentResponse{ID: result.ID.String()}, nil
}
