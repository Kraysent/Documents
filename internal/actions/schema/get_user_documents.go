package schema

import (
	"documents/internal/library/validation"
)

type GetUserDocumentsRequest struct {
	Username string `schema:"username"`
}

func (r *GetUserDocumentsRequest) Validate() error {
	rules := validation.NewRules(
		validation.StringNotEmpty(r.Username),
	)

	return rules.Validate()
}

type GetUserDocumentsResponse struct {
	Documents []GetDocumentResponse `json:"documents"`
}
