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

	data, err := repo.Storages.Documents.GetDocument(ctx, documents.GetDocumentRequest{ID: id})
	if err == sql.ErrNoRows {
		return nil, web.NotFoundError(fmt.Errorf("active user does not have document with ID %s", r.ID))
	}
	if err != nil {
		return nil, web.DatabaseError(err)
	}

	if data.Document.Owner != userID {
		return nil, web.NotFoundError(fmt.Errorf("active user does not have document with ID %s", r.ID))
	}

	return &schema.GetDocumentResponse{
		ID:          data.Document.ID.String(),
		Name:        data.Document.Name,
		Version:     data.Document.Version,
		Description: data.Document.Description,
	}, nil
}
