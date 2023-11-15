package actions

import (
	"context"
	"fmt"

	"documents/internal/actions/schema"
	"documents/internal/core"
	"documents/internal/library/web"
	"documents/internal/storage/documents"
)

func GetUserDocuments(
	ctx context.Context, repo *core.Repository, r schema.GetUserDocumentsRequest,
) (*schema.GetUserDocumentsResponse, error) {
	userID := repo.SessionManager.GetInt64(ctx, "user_id")
	if userID == 0 {
		return nil, web.AuthorizationError(fmt.Errorf("failed to authorize"))
	}

	result, err := repo.Storages.Documents.GetDocuments(ctx, documents.GetDocumentsRequest{
		Fields: map[string]any{
			documents.ColumnOwner: userID,
		},
	})
	if err != nil {
		return nil, web.DatabaseError(err)
	}

	response := schema.GetUserDocumentsResponse{
		Documents: []schema.GetDocumentResponse{},
	}

	for _, doc := range result.Documents {
		response.Documents = append(response.Documents, schema.GetDocumentResponse{
			ID:          doc.ID.String(),
			Name:        doc.Name,
			Version:     doc.Version,
			Description: doc.Description,
		})
	}

	return &response, nil
}
