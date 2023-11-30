package schema

import (
	"documents/internal/library/validation"
)

type CreateLinkRequest struct {
	UserID     int64  `json:"-"`
	DocumentID string `json:"document_id"`
	ExpiryDate string `json:"expiry_date"`
}

func (r *CreateLinkRequest) Validate() error {
	rules := validation.NewRules(
		validation.IsUUID(r.DocumentID),
		validation.IsISO8601(r.ExpiryDate),
	)

	return rules.Validate()
}

type CreateLinkResponse struct {
	ID string `json:"id"`
}
