package actions

import (
	"context"

	"documents/internal/core"
	"documents/internal/entities"
)

func InsertDocument(ctx context.Context, repo *core.Repository, data entities.Document) (string, *entities.CodedError) {
	id, err := repo.Storage.DocumentStorage.AddDocument(ctx, data.Username, data.Type, data.Attributes)
	if err != nil {
		return "", entities.DatabaseError(err)
	}

	return id, nil
}
