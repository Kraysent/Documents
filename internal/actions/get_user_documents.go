package actions

import (
	"context"
	"encoding/hex"

	"documents/internal/actions/schema"
	"documents/internal/core"
	"documents/internal/storage"
)

func GetUserDocuments(
	ctx context.Context, repo *core.Repository, request schema.GetUserDocumentsRequest,
) (*schema.GetUserDocumentsResponse, error) {
	result, err := repo.Storage.DocumentStorage.GetDocuments(ctx, storage.GetDocumentsRequest{
		Fields: map[string]any{
			storage.DocumentsColumnUsername: request.Username,
		},
	})
	if err != nil {
		return nil, err
	}

	var response schema.GetUserDocumentsResponse

	for _, doc := range result.Documents {
		response.Documents = append(response.Documents, schema.GetDocumentResponse{
			ID:           hex.EncodeToString(doc.ID),
			DocumentType: doc.Type,
			Attributes:   doc.Attributes,
		})
	}

	return &response, err
}
