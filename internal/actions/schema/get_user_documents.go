package schema

type GetUserDocumentsRequest struct{}

func (r *GetUserDocumentsRequest) Validate() error {
	return nil
}

type GetUserDocumentsResponse struct {
	Documents []GetDocumentResponse `json:"documents"`
}
