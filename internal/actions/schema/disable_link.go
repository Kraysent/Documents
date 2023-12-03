package schema

import (
	"documents/internal/library/validation"
)

type DisableLinkRequest struct {
	ID string `schema:"id"`
}

func (r DisableLinkRequest) Validate() error {
	rules := validation.NewRules(
		validation.StringNotEmpty(r.ID),
		validation.IsUUID(r.ID),
	)

	return rules.Validate()
}

type DisableLinkResponse struct {
	Status string `json:"status"`
}
