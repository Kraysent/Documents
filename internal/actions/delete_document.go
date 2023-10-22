package actions

import (
	"context"
	"encoding/hex"
	"fmt"

	"documents/internal/actions/schema"
	"documents/internal/core"
	"documents/internal/library/web"
	"documents/internal/storage/documents"
)

func DeleteDocument(
	ctx context.Context, repo *core.Repository, r schema.DeleteDocumentRequest,
) (*schema.DeleteDocumentResponse, error) {
	userID := repo.SessionManager.GetInt64(ctx, "user_id")
	if userID == 0 {
		return nil, web.AuthorizationError(fmt.Errorf("failed to authorize"))
	}

	id, err := hex.DecodeString(r.ID)
	if err != nil {
		return nil, web.ValidationError(err)
	}

	n, err := repo.Storages.Documents.RemoveDocuments(ctx,
		map[string]any{
			documents.ColumnID: id,
		},
	)
	if err != nil {
		return nil, web.DatabaseError(err)
	}

	return &schema.DeleteDocumentResponse{DeletedNumber: n}, nil
}
