package actions

import (
	"context"
	"fmt"

	"documents/internal/core"
	"documents/internal/entities"
	"documents/internal/storage"
)

func GetDocumentByID(ctx context.Context, repo *core.Repository, id string) (*entities.Document, *entities.CodedError) {
	data, err := repo.Storage.DocumentStorage.GetDocuments(ctx,
		storage.GetDocumentsRequest{Fields: map[string]any{"id": id}},
	)
	if err != nil {
		return nil, entities.DatabaseError(err)
	}

	if len(data.Documents) != 1 {
		return nil, entities.InternalError(fmt.Errorf("database returned %d rows, expected 1", len(data.Documents)))
	}

	return &data.Documents[0], nil
}

func GetDocumentByUsernameAndType(
	ctx context.Context, repo *core.Repository, username string, documentType string,
) (*entities.Document, *entities.CodedError) {
	data, err := repo.Storage.DocumentStorage.GetDocuments(
		ctx, storage.GetDocumentsRequest{
			Fields: map[string]any{"username": username, "document_type": documentType},
		},
	)
	if err != nil {
		return nil, entities.DatabaseError(err)
	}

	if len(data.Documents) != 1 {
		return nil, entities.InternalError(fmt.Errorf("database returned %d rows, expected 1", len(data.Documents)))
	}

	return &data.Documents[0], nil
}
