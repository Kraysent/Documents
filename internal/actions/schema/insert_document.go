package schema

import (
	"documents/internal/library/validation"
)

type InsertDocumentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r InsertDocumentRequest) Validate() error {
	rules := validation.NewRules(
		validation.StringNotEmpty(r.Name),
	)

	return rules.Validate()
}

type InsertDocumentResponse struct {
	ID string `json:"id"`
}
