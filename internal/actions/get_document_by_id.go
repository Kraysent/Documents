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

func GetDocumentByID(
	ctx context.Context, repo *core.Repository, r schema.GetDocumentByIDRequest,
) (*schema.GetDocumentResponse, error) {
	userID := repo.SessionManager.GetInt64(ctx, "user_id")
	if userID == 0 {
		return nil, web.AuthorizationError(fmt.Errorf("failed to authorize"))
	}

	idBytes, err := hex.DecodeString(r.ID)
	if err != nil {
		return nil, web.ValidationError(err)
	}

	data, err := repo.Storages.Documents.GetDocuments(ctx,
		documents.GetDocumentsRequest{Fields: map[string]any{documents.ColumnID: idBytes}},
	)
	if err != nil {
		return nil, web.DatabaseError(err)
	}

	if len(data.Documents) != 1 {
		return nil, web.InternalError(fmt.Errorf("database returned %d rows, expected 1", len(data.Documents)))
	}

	if data.Documents[0].UserID != userID {
		return nil, web.AuthorizationError(fmt.Errorf("active user does not have document with this ID"))
	}

	return &schema.GetDocumentResponse{
		ID:           hex.EncodeToString(data.Documents[0].ID),
		DocumentType: data.Documents[0].Type,
		Attributes:   data.Documents[0].Attributes,
	}, nil
}
