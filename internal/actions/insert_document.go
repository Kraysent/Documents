package actions

import (
	"context"

	"documents/internal/core"
	"documents/internal/entities"
)

func InsertDocument(ctx context.Context, repo *core.Repository, data entities.Document) (string, *entities.CodedError) {
	dbData := make(map[string]any)
	dbData["username"] = data.Username
	dbData["document_type"] = data.Type

	for key, value := range data.Attributes {
		dbData[key] = value
	}

	id, err := repo.Storage.DocumentStorage.AddDocument(ctx, dbData)
	if err != nil {
		return "", entities.DatabaseError(err)
	}

	return id, nil
}
