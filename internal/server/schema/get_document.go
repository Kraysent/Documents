package schema

import (
	"fmt"
)

type GetDocumentByIDRequest struct {
	ID string `schema:"id"`
}

func (r *GetDocumentByIDRequest) Validate() error {
	if r.ID == "" {
		return fmt.Errorf("empty id")
	}

	return nil
}
