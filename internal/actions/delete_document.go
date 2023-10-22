package actions

import (
	"context"
	"encoding/hex"

	"documents/internal/actions/schema"
	"documents/internal/core"
	"documents/internal/library/web"
	"documents/internal/storage/documents"
)

func DeleteDocument(
	ctx context.Context, repo *core.Repository, r schema.DeleteDocumentRequest,
) (*schema.DeleteDocumentResponse, error) {
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
