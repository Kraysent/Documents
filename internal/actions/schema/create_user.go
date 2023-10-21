package schema

import (
	"documents/internal/validation"
)

type CreateUserRequest struct {
	Username string `json:"username"`
}

func (r *CreateUserRequest) Validate() error {
	rules := validation.NewRules(
		validation.StringNotEmpty(r.Username),
	)

	return rules.Validate()
}

type CreateUserResponse struct {
	Status string `json:"status"`
}
