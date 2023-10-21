package schema

import (
	"documents/internal/validation"
)

type GetDocumentByIDRequest struct {
	ID string `schema:"id"`
}

func (r *GetDocumentByIDRequest) Validate() error {
	rules := validation.NewRules(
		validation.StringNotEmpty(r.ID),
	)

	return rules.Validate()
}

type GetDocumentResponse struct {
	ID           string         `json:"id"`
	DocumentType string         `json:"document_type"`
	Attributes   map[string]any `json:"attributes"`
}

type GetDocumentByUsernameAndTypeRequest struct {
	Username string `schema:"username"`
	Type     string `schema:"document_type"`
}

func (r *GetDocumentByUsernameAndTypeRequest) Validate() error {
	rules := validation.NewRules(
		validation.StringNotEmpty(r.Username),
		validation.StringNotEmpty(r.Type),
	)

	return rules.Validate()
}

type GetDocumentsResponse struct {
	Documents []GetDocumentResponse `json:"documents"`
}
