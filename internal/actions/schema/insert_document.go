package schema

import (
	"documents/internal/library/validation"
)

type InsertDocumentRequest struct {
	Username   string         `json:"username" db:"username"`
	Type       string         `json:"document_type" db:"document_type"`
	Attributes map[string]any `json:"attributes" db:"attributes"`
}

func (r *InsertDocumentRequest) Validate() error {
	rules := validation.NewRules(
		validation.StringNotEmpty(r.Username),
		validation.StringNotEmpty(r.Type),
	)

	return rules.Validate()
}

type InsertDocumentResponse struct {
	ID string `json:"id"`
}
