package actions

import (
	"context"
	"fmt"

	"documents/internal/actions/schema"
	"documents/internal/core"
	"documents/internal/library/web"
	"documents/internal/storage/documents"
	"github.com/google/uuid"
)

func GetDocumentByID(
	ctx context.Context, repo *core.Repository, r schema.GetDocumentByIDRequest,
) (*schema.GetDocumentResponse, error) {
	userID := repo.SessionManager.GetInt64(ctx, "user_id")
	if userID == 0 {
		return nil, web.AuthorizationError(fmt.Errorf("failed to authorize"))
	}

	id, err := uuid.Parse(r.ID)
	if err != nil {
		return nil, web.ValidationError(err)
	}

	data, err := repo.Storages.Documents.GetDocuments(ctx,
		documents.GetDocumentsRequest{Fields: map[string]any{documents.ColumnID: id}},
	)
	if err != nil {
		return nil, web.DatabaseError(err)
	}

	if len(data.Documents) != 1 {
		return nil, web.InternalError(fmt.Errorf("database returned %d rows, expected 1", len(data.Documents)))
	}

	if data.Documents[0].Owner != userID {
		return nil, web.AuthorizationError(fmt.Errorf("active user does not have document with this ID"))
	}

	return &schema.GetDocumentResponse{
		ID:          data.Documents[0].ID.String(),
		Name:        data.Documents[0].Name,
		Version:     data.Documents[0].Version,
		Description: data.Documents[0].Description,
	}, nil
}
