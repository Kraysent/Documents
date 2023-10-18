package schema

import (
	"fmt"
)

type DeleteDocumentRequest struct {
	Username string `schema:"username"`
	Type     string `schema:"document_type"`
}

func (r *DeleteDocumentRequest) Validate() error {
	if r.Username == "" {
		return fmt.Errorf("empty username")
	}

	if r.Type == "" {
		return fmt.Errorf("empty document_type")
	}

	return nil
}

type DeleteDocumentResponse struct {
	DeletedNumber int64 `json:"deleted_number"`
}
