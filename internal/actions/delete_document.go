package actions

import (
	"context"
	"database/sql"
	"fmt"

	"documents/internal/actions/schema"
	"documents/internal/core"
	"documents/internal/library/web"
	"documents/internal/storage/documents"
	"github.com/google/uuid"
)

func DeleteDocument(
	ctx context.Context, repo *core.Repository, r schema.DeleteDocumentRequest,
) (*schema.DeleteDocumentResponse, error) {
	userID := repo.SessionManager.GetInt64(ctx, "user_id")
	if userID == 0 {
		return nil, web.AuthorizationError(fmt.Errorf("failed to authorize"))
	}

	id, err := uuid.Parse(r.ID)
	if err != nil {
		return nil, web.ValidationError(err)
	}

	res, err := repo.Storages.Documents.GetDocument(ctx, documents.GetDocumentRequest{ID: id})
	if err == sql.ErrNoRows {
		return nil, web.NotFoundError(fmt.Errorf("active user does not have document with ID %s", r.ID))
	}
	if res.Document.Owner != userID {
		return nil, web.NotFoundError(fmt.Errorf("active user does not have document with ID %s", r.ID))
	}

	n, err := repo.Storages.Documents.RemoveDocuments(ctx,
		map[string]any{
			documents.ColumnID: id,
		},
	)
	if err != nil {
		return nil, web.DatabaseError(err)
	}
	if n == 0 {
		return nil, web.InternalError(fmt.Errorf("unable to remove document with ID %s", r.ID))
	}

	return &schema.DeleteDocumentResponse{DeletedNumber: n}, nil
}
