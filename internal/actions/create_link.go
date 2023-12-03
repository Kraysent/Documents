package actions

import (
	"context"
	"database/sql"
	"errors"
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
	expiryDate, err := time.Parse(time.RFC3339, r.ExpiryDate)
	if err != nil {
		return nil, web.ValidationError(err)
	}

	if expiryDate.Before(time.Now()) {
		return nil, web.ValidationError(fmt.Errorf("expiry date must be after current time"))
	}

	documentID, err := uuid.Parse(r.DocumentID)
	if err != nil {
		return nil, web.ValidationError(err)
	}

	data, err := repo.Storages.Documents.GetDocument(ctx,
		documents.GetDocumentRequest{ID: documentID},
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, web.NotFoundError(fmt.Errorf("active user does not have document with ID %s", r.DocumentID))
	}
	if err != nil {
		return nil, web.DatabaseError(err)
	}

	if data.Document.Owner != r.UserID {
		return nil, web.NotFoundError(fmt.Errorf("active user does not have document with ID %s", r.DocumentID))
	}

	result, err := repo.Storages.Links.CreateLink(ctx, links.CreateLinkRequest{
		DocumentID: documentID,
		ExpiryDate: expiryDate,
	})

	return &schema.CreateLinkResponse{ID: result.ID.String()}, nil
}
