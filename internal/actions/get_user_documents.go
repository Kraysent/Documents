package actions

import (
	"context"
	"encoding/hex"

	"documents/internal/actions/schema"
	"documents/internal/core"
	"documents/internal/library/web"
	"documents/internal/storage/documents"
	"documents/internal/storage/users"
)

func GetUserDocuments(
	ctx context.Context, repo *core.Repository, r schema.GetUserDocumentsRequest,
) (*schema.GetUserDocumentsResponse, error) {
	res, err := repo.Storages.Users.GetUser(ctx, users.GetUserRequest{Fields: map[string]any{
		users.ColumnUsername: r.Username,
	}})
	if err != nil {
		return nil, web.DatabaseError(err)
	}

	result, err := repo.Storages.Documents.GetDocuments(ctx, documents.GetDocumentsRequest{
		Fields: map[string]any{
			documents.ColumnUserID: res.UserID,
		},
	})
	if err != nil {
		return nil, web.DatabaseError(err)
	}

	var response schema.GetUserDocumentsResponse

	for _, doc := range result.Documents {
		response.Documents = append(response.Documents, schema.GetDocumentResponse{
			ID:           hex.EncodeToString(doc.ID),
			DocumentType: doc.Type,
			Attributes:   doc.Attributes,
		})
	}

	return &response, nil
}
