package schema

import (
	"documents/internal/library/validation"
)

type GetDocumentByLinkRequest struct {
	ID string `schema:"id"`
}

func (r GetDocumentByLinkRequest) Validate() error {
	rules := validation.NewRules(
		validation.StringNotEmpty(r.ID),
		validation.IsUUID(r.ID),
	)

	return rules.Validate()
}
