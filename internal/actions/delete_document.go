package actions

import (
	"context"

	"documents/internal/actions/schema"
	"documents/internal/core"
	"documents/internal/entities"
	"documents/internal/storage"
)

func DeleteDocument(
	ctx context.Context, repo *core.Repository, request schema.DeleteDocumentRequest,
) (*schema.DeleteDocumentResponse, *entities.CodedError) {
	n, err := repo.Storage.DocumentStorage.RemoveDocument(ctx,
		map[string]any{
			storage.DocumentsColumnUsername:     request.Username,
			storage.DocumentsColumnDocumentType: request.Type,
		},
	)
	if err != nil {
		return nil, entities.DatabaseError(err)
	}

	return &schema.DeleteDocumentResponse{DeletedNumber: n}, nil
}
