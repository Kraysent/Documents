package schema

import (
	"documents/internal/library/validation"
)

type GetOrCreateUserRequest struct {
	GoogleID string `json:"id"`
	Email    string `json:"email"`
}

func (r *GetOrCreateUserRequest) Validate() error {
	rules := validation.NewRules(
		validation.StringNotEmpty(r.GoogleID),
		validation.StringNotEmpty(r.Email),
	)

	return rules.Validate()
}

type GetOrCreateUserResponse struct {
	Status string `json:"status"`
}
