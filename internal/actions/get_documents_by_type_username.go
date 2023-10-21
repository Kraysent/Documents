package actions

import (
	"context"
	"encoding/hex"

	"documents/internal/actions/schema"
	"documents/internal/core"
	"documents/internal/storage/documents"
	"documents/internal/storage/users"
)

func GetDocumentsByUsernameAndType(
	ctx context.Context, repo *core.Repository, r schema.GetDocumentByUsernameAndTypeRequest,
) (*schema.GetDocumentsResponse, error) {
	res, err := repo.Storages.Users.GetUser(ctx, users.GetUserRequest{Username: r.Username})
	if err != nil {
		return nil, DatabaseError(err)
	}

	data, err := repo.Storages.Documents.GetDocuments(
		ctx, documents.GetDocumentsRequest{
			Fields: map[string]any{
				documents.ColumnUserID:       res.UserID,
				documents.ColumnDocumentType: r.Type,
			},
		},
	)
	if err != nil {
		return nil, DatabaseError(err)
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
