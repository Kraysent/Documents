package actions

import (
	"context"
	"encoding/hex"
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
			documents.ColumnUserID: userID,
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
			ID:           hex.EncodeToString(doc.ID),
			DocumentType: doc.Type,
			Attributes:   doc.Attributes,
		})
	}

	return &response, nil
}
