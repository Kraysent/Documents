package documents

import (
	"github.com/google/uuid"
)

const (
	TableName         = "documents.t_document"
	ColumnID          = "id"
	ColumnName        = "name"
	ColumnOwner       = "owner"
	ColumnVersion     = "version"
	ColumnDescription = "description"
)

type Document struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Owner       int64     `db:"owner"`
	Version     int64     `db:"version"`
	Description string    `db:"description"`
}

type AddDocumentRequest struct {
	Name        string
	Owner       int64
	Description string
}

type AddDocumentResult struct {
	ID uuid.UUID
}

type GetDocumentsRequest struct {
	Fields map[string]any
}

type GetDocumentsResult struct {
	Documents []Document
}
