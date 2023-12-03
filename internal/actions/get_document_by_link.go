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

func GetDocumentByLink(
	ctx context.Context, repo *core.Repository, r schema.GetDocumentByLinkRequest,
) (*schema.GetDocumentResponse, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return nil, web.ValidationError(err)
	}

	result, err := repo.Storages.Links.GetLink(ctx, links.GetLinkRequest{ID: id})
	if err != nil {
		return nil, web.DatabaseError(err)
	}

	if result.Link.Status != links.StatusEnabled {
		return nil, web.NotFoundError(fmt.Errorf("link is unavailible"))
	}

	if result.Link.ExpiryDate.Before(time.Now()) {
		return nil, web.NotFoundError(fmt.Errorf("link is expired"))
	}

	docResult, err := repo.Storages.Documents.GetDocument(ctx, documents.GetDocumentRequest{
		ID: result.Link.DocumentID,
	})
	if err != nil {
		return nil, web.DatabaseError(err)
	}

	return &schema.GetDocumentResponse{
		ID:          docResult.Document.ID.String(),
		Name:        docResult.Document.Name,
		Version:     docResult.Document.Version,
		Description: docResult.Document.Description,
	}, nil
}
