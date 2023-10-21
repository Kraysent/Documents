package actions

import (
	"context"
	"encoding/hex"
	"fmt"

	"documents/internal/actions/schema"
	"documents/internal/core"
	"documents/internal/storage/documents"
)

func GetDocumentByID(
	ctx context.Context, repo *core.Repository, r schema.GetDocumentByIDRequest,
) (*schema.GetDocumentResponse, error) {
	idBytes, err := hex.DecodeString(r.ID)
	if err != nil {
		return nil, ValidationError(err)
	}

	data, err := repo.Storages.Documents.GetDocuments(ctx,
		documents.GetDocumentsRequest{Fields: map[string]any{documents.ColumnID: idBytes}},
	)
	if err != nil {
		return nil, DatabaseError(err)
	}

	if len(data.Documents) != 1 {
		return nil, InternalError(fmt.Errorf("database returned %d rows, expected 1", len(data.Documents)))
	}

	return &schema.GetDocumentResponse{
		ID:           hex.EncodeToString(data.Documents[0].ID),
		DocumentType: data.Documents[0].Type,
		Attributes:   data.Documents[0].Attributes,
	}, nil
}
