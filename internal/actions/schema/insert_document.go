package schema

import (
	"documents/internal/library/validation"
)

type InsertDocumentRequest struct {
	Type       string         `json:"document_type" db:"document_type"`
	Attributes map[string]any `json:"attributes" db:"attributes"`
}

func (r *InsertDocumentRequest) Validate() error {
	rules := validation.NewRules(
		validation.StringNotEmpty(r.Type),
	)

	return rules.Validate()
}

type InsertDocumentResponse struct {
	ID string `json:"id"`
}
