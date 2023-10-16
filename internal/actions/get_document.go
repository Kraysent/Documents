package actions

import (
	"context"
	"fmt"

	"documents/internal/core"
	"documents/internal/entities"
)

func GetDocumentByID(ctx context.Context, repo *core.Repository, id string) (*entities.Document, *entities.CodedError) {
	data, err := repo.Storage.DocumentStorage.GetDocumentByID(ctx, id)
	if err != nil {
		return nil, entities.DatabaseError(err)
	}

	var document entities.Document

	usernameAny, ok := data["username"]
	if !ok {
		return nil, entities.InternalError(fmt.Errorf("cannot get username from database response"))
	}
	username, ok := usernameAny.(string)
	if !ok {
		return nil, entities.InternalError(fmt.Errorf("username from database response has wrong type"))
	}

	document.Username = username
	delete(data, "username")

	documentTypeAny, ok := data["document_type"]
	if !ok {
		return nil, entities.InternalError(fmt.Errorf("cannot get document type from database response"))
	}
	documentType, ok := documentTypeAny.(string)
	if !ok {
		return nil, entities.InternalError(fmt.Errorf("document type from database response has wrong type"))
	}

	document.Type = documentType
	delete(data, "document_type")

	document.Attributes = data

	return &document, nil
}
