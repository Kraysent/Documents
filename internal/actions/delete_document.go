package actions

import (
	"context"

	"documents/internal/core"
	"documents/internal/entities"
)

func DeleteDocument(
	ctx context.Context, repo *core.Repository, username string, documentType string,
) (map[string]any, *entities.CodedError) {
	n, err := repo.Storage.DocumentStorage.RemoveDocument(ctx, map[string]any{"username": username, "document_type": documentType})
	if err != nil {
		return nil, entities.DatabaseError(err)
	}

	return map[string]any{"deleted_number": n}, nil
}
