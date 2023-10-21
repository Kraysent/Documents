package actions

import (
	"context"

	"documents/internal/actions/schema"
	"documents/internal/core"
	"documents/internal/storage/users"
)

func InsertDocument(
	ctx context.Context, repo *core.Repository, r schema.InsertDocumentRequest,
) (*schema.InsertDocumentResponse, error) {
	res, err := repo.Storages.Users.GetUser(ctx, users.GetUserRequest{Username: r.Username})
	if err != nil {
		return nil, DatabaseError(err)
	}

	id, err := repo.Storages.Documents.AddDocument(ctx, res.UserID, r.Type, r.Attributes)
	if err != nil {
		return nil, DatabaseError(err)
	}

	return &schema.InsertDocumentResponse{ID: id}, nil
}
