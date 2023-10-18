package actions

import (
	"context"
	"encoding/hex"
	"fmt"

	"documents/internal/actions/schema"
	"documents/internal/core"
	"documents/internal/entities"
	"documents/internal/storage"
)

func GetDocumentByID(
	ctx context.Context, repo *core.Repository, request schema.GetDocumentByIDRequest,
) (*schema.GetDocumentResponse, *entities.CodedError) {
	idBytes, err := hex.DecodeString(request.ID)
	if err != nil {
		return nil, entities.ValidationError(err)
	}

	data, err := repo.Storage.DocumentStorage.GetDocuments(ctx,
		storage.GetDocumentsRequest{Fields: map[string]any{storage.DocumentsColumnID: idBytes}},
	)
	if err != nil {
		return nil, entities.DatabaseError(err)
	}

	if len(data.Documents) != 1 {
		return nil, entities.InternalError(fmt.Errorf("database returned %d rows, expected 1", len(data.Documents)))
	}

	return &schema.GetDocumentResponse{
		ID:           hex.EncodeToString(data.Documents[0].ID),
		DocumentType: data.Documents[0].Type,
		Attributes:   data.Documents[0].Attributes,
	}, nil
}

func GetDocumentByUsernameAndType(
	ctx context.Context, repo *core.Repository, request schema.GetDocumentByUsernameAndTypeRequest,
) (*schema.GetDocumentResponse, *entities.CodedError) {
	data, err := repo.Storage.DocumentStorage.GetDocuments(
		ctx, storage.GetDocumentsRequest{
			Fields: map[string]any{
				storage.DocumentsColumnUsername:     request.Username,
				storage.DocumentsColumnDocumentType: request.Type,
			},
		},
	)
	if err != nil {
		return nil, entities.DatabaseError(err)
	}

	if len(data.Documents) != 1 {
		return nil, entities.InternalError(fmt.Errorf("database returned %d rows, expected 1", len(data.Documents)))
	}

	return &schema.GetDocumentResponse{
		ID:           hex.EncodeToString(data.Documents[0].ID),
		DocumentType: data.Documents[0].Type,
		Attributes:   data.Documents[0].Attributes,
	}, nil
}
