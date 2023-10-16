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

type GetDocumentByUsernameAndTypeRequest struct {
	Username string `schema:"username"`
	Type     string `schema:"document_type"`
}

func (r *GetDocumentByUsernameAndTypeRequest) Validate() error {
	if r.Username == "" {
		return fmt.Errorf("empty username")
	}

	if r.Type == "" {
		return fmt.Errorf("empty document_type")
	}

	return nil
}
