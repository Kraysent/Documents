package storage

import (
	"documents/internal/entities"
)

type GetDocumentsRequest struct {
	Fields map[string]any
}

type GetDocumentsResult struct {
	Documents []entities.Document
}
