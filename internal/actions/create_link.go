package actions

import (
	"context"
	"fmt"
	"time"

	"documents/internal/actions/schema"
	"documents/internal/core"
	"documents/internal/library/web"
	"documents/internal/storage/documents"
	"documents/internal/storage/links"
	"github.com/google/uuid"
)

func CreateLink(
	ctx context.Context, repo *core.Repository, r schema.CreateLinkRequest,
) (*schema.CreateLinkResponse, error) {
	documentID, err := uuid.Parse(r.DocumentID)
	if err != nil {
		return nil, web.ValidationError(err)
	}

	data, err := repo.Storages.Documents.GetDocuments(ctx,
		documents.GetDocumentsRequest{Fields: map[string]any{documents.ColumnID: documentID}},
	)
	if err != nil {
		return nil, web.DatabaseError(err)
	}

	if len(data.Documents) != 1 {
		return nil, web.InternalError(fmt.Errorf("database returned %d rows, expected 1", len(data.Documents)))
	}

	if data.Documents[0].Owner != r.UserID {
		return nil, web.AuthorizationError(fmt.Errorf("user does not have document with this ID"))
	}

	expiryDate, err := time.Parse(time.RFC3339, r.ExpiryDate)
	if err != nil {
		return nil, web.ValidationError(err)
	}

	result, err := repo.Storages.Links.CreateLink(ctx, links.CreateLinkRequest{
		DocumentID: documentID,
		ExpiryDate: expiryDate,
	})

	return &schema.CreateLinkResponse{ID: result.ID.String()}, nil
}
