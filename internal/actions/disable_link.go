package actions

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"documents/internal/actions/schema"
	"documents/internal/core"
	"documents/internal/library/web"
	"documents/internal/storage/documents"
	"documents/internal/storage/links"
	"github.com/google/uuid"
)

func DisableLink(
	ctx context.Context, repo *core.Repository, r schema.DisableLinkRequest,
) (*schema.DisableLinkResponse, error) {
	userID := repo.SessionManager.GetInt64(ctx, "user_id")
	if userID == 0 {
		return nil, web.AuthorizationError(fmt.Errorf("failed to authorize"))
	}

	linkID, err := uuid.Parse(r.ID)
	if err != nil {
		return nil, web.ValidationError(err)
	}

	linkResult, err := repo.Storages.Links.GetLink(ctx, links.GetLinkRequest{ID: linkID})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, web.NotFoundError(fmt.Errorf("active user does not have link with ID %s", r.ID))
	}
	if err != nil {
		return nil, web.DatabaseError(err)
	}

	docResult, err := repo.Storages.Documents.GetDocument(ctx, documents.GetDocumentRequest{
		ID: linkResult.Link.DocumentID,
	})
	if err != nil {
		return nil, web.DatabaseError(err)
	}

	if docResult.Document.Owner != userID {
		return nil, web.NotFoundError(fmt.Errorf("active user does not have link with ID %s", r.ID))
	}
	if linkResult.Link.Status != links.StatusEnabled {
		return nil, web.ValidationError(fmt.Errorf("link already disabled"))
	}

	_, err = repo.Storages.Links.SetLinkStatus(ctx, links.SetLinkStatusRequest{
		ID:     linkID,
		Status: links.StatusDisabled,
	})
	if err != nil {
		return nil, web.DatabaseError(err)
	}

	return &schema.DisableLinkResponse{Status: string(links.StatusDisabled)}, nil
}
