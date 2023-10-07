package actions

import (
	"context"
	"fmt"

	"documents/internal/core"
)

func InsertDocument(ctx context.Context, repo *core.Repository) error {
	id, err := repo.Storage.DocumentStorage.AddDocument(ctx, map[string]any{"test": "lala1", "type": "passport"})
	if err != nil {
		return err
	}

	fmt.Println(id)

	return nil
}
