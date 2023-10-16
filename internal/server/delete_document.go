package server

import (
	"context"
	"net/http"

	"documents/internal/actions"
	"documents/internal/entities"
	schema2 "documents/internal/server/schema"

	"github.com/gorilla/schema"
)

func (s *Server) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var request schema2.DeleteDocumentRequest
	decoder := schema.NewDecoder()

	if err := decoder.Decode(&request, r.URL.Query()); err != nil {
		s.handleError(w, entities.ValidationError(err))
		return
	}

	if err := request.Validate(); err != nil {
		s.handleError(w, entities.ValidationError(err))
		return
	}

	document, cErr := actions.DeleteDocument(ctx, s.repo, request.Username, request.Type)
	if cErr != nil {
		s.handleError(w, cErr)
		return
	}

	if err := s.handleOK(w, document); err != nil {
		s.handleError(w, entities.ValidationError(err))
		return
	}
}
