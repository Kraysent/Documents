package schema

import (
	"documents/internal/validation"
)

type DeleteDocumentRequest struct {
	ID string `schema:"id"`
}

func (r *DeleteDocumentRequest) Validate() error {
	rules := validation.NewRules(
		validation.StringNotEmpty(r.ID),
	)

	return rules.Validate()
}

type DeleteDocumentResponse struct {
	DeletedNumber int64 `json:"deleted_number"`
}
