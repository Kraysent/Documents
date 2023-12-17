package schema

import (
	"documents/internal/library/validation"
)

type GetDocumentByIDRequest struct {
	ID string `schema:"id"`
}

func (r GetDocumentByIDRequest) Validate() error {
	rules := validation.NewRules(
		validation.IsUUID(r.ID),
	)

	return rules.Validate()
}

type GetDocumentResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     int64  `json:"version"`
	Description string `json:"description"`
}

type GetDocumentsResponse struct {
	Documents []GetDocumentResponse `json:"documents"`
}
