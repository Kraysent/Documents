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

const (
	DefaultPageSize = 25
)

func GetLinks(
	ctx context.Context, repo *core.Repository, r schema.GetLinksRequest,
) (*schema.GetLinksResponse, error) {
	userID := repo.SessionManager.GetInt64(ctx, "user_id")
	if userID == 0 {
		return nil, web.AuthorizationError(fmt.Errorf("failed to authorize"))
	}

	if r.PageSize == 0 {
		r.PageSize = DefaultPageSize
	}

	documentID, err := uuid.Parse(r.DocumentID)
	if err != nil {
		return nil, web.ValidationError(err)
	}

	docResult, err := repo.Storages.Documents.GetDocument(ctx, documents.GetDocumentRequest{
		ID: documentID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, web.NotFoundError(fmt.Errorf("active user does not have document with ID %s", r.DocumentID))
	}
	if err != nil {
		return nil, web.DatabaseError(err)
	}

	if docResult.Document.Owner != userID {
		return nil, web.NotFoundError(fmt.Errorf("active user does not have document with ID %s", r.DocumentID))
	}

	fields := map[string]any{
		links.ColumnDocumentID: r.DocumentID,
	}
	if r.Status != "" {
		fields[links.ColumnStatus] = r.Status
	}

	linksResult, err := repo.Storages.Links.GetLinks(ctx, links.GetLinksRequest{
		Fields:       fields,
		PageSize:     r.PageSize,
		PageNumber:   r.Page,
		OrderByField: links.ColumnExpiryDate,
	})

	var result schema.GetLinksResponse
	result.Links = make([]schema.Link, len(linksResult.Links))

	for i, link := range linksResult.Links {
		result.Links[i] = schema.Link{
			ID:           link.ID.String(),
			DocumentID:   link.DocumentID.String(),
			CreationDate: link.CreationDate.Format(time.RFC3339),
			ExpiryDate:   link.ExpiryDate.Format(time.RFC3339),
			Status:       string(link.Status),
		}
	}

	return &result, nil
}
