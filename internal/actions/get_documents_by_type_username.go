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

func GetDocumentsByUsernameAndType(
	ctx context.Context, repo *core.Repository, r schema.GetDocumentByUsernameAndTypeRequest,
) (*schema.GetDocumentsResponse, error) {
	userID := repo.SessionManager.GetInt64(ctx, "user_id")
	if userID == 0 {
		return nil, web.AuthorizationError(fmt.Errorf("failed to authorize"))
	}

	data, err := repo.Storages.Documents.GetDocuments(
		ctx, documents.GetDocumentsRequest{
			Fields: map[string]any{
				documents.ColumnUserID:       userID,
				documents.ColumnDocumentType: r.Type,
			},
		},
	)
	if err != nil {
		return nil, web.DatabaseError(err)
	}

	var response schema.GetDocumentsResponse

	for _, doc := range data.Documents {
		response.Documents = append(response.Documents, schema.GetDocumentResponse{
			ID:           hex.EncodeToString(doc.ID),
			DocumentType: doc.Type,
			Attributes:   doc.Attributes,
		})
	}

	return &response, nil
}
