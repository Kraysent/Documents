package actions

import (
	"context"
	"fmt"

	"documents/internal/core"
	"documents/internal/entities"
)

func attributesToDocument(attributes map[string]any) (*entities.Document, error) {
	var document entities.Document

	usernameAny, ok := attributes["username"]
	if !ok {
		return nil, fmt.Errorf("cannot get username from database response")
	}
	username, ok := usernameAny.(string)
	if !ok {
		return nil, fmt.Errorf("username from database response has wrong type")
	}

	document.Username = username
	delete(attributes, "username")

	documentTypeAny, ok := attributes["document_type"]
	if !ok {
		return nil, fmt.Errorf("cannot get document type from database response")
	}
	documentType, ok := documentTypeAny.(string)
	if !ok {
		return nil, fmt.Errorf("document type from database response has wrong type")
	}

	document.Type = documentType
	delete(attributes, "document_type")

	document.Attributes = attributes
	return &document, nil
}

func GetDocumentByID(ctx context.Context, repo *core.Repository, id string) (*entities.Document, *entities.CodedError) {
	data, err := repo.Storage.DocumentStorage.GetDocument(ctx, map[string]any{"id": id})
	if err != nil {
		return nil, entities.DatabaseError(err)
	}

	document, err := attributesToDocument(data)
	if err != nil {
		return nil, entities.InternalError(err)
	}

	return document, nil
}

func GetDocumentByUsernameAndType(
	ctx context.Context, repo *core.Repository, username string, documentType string,
) (*entities.Document, *entities.CodedError) {
	data, err := repo.Storage.DocumentStorage.GetDocument(
		ctx, map[string]any{"username": username, "document_type": documentType},
	)
	if err != nil {
		return nil, entities.DatabaseError(err)
	}

	document, err := attributesToDocument(data)
	if err != nil {
		return nil, entities.InternalError(err)
	}

	return document, nil
}
