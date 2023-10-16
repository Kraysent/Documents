package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"documents/internal/actions"
	"documents/internal/entities"
)

func (s *Server) InsertDocument(writer http.ResponseWriter, request *http.Request) {
	ctx := context.Background()
	var document entities.Document

	data, err := io.ReadAll(request.Body)
	if err != nil {
		s.handleError(writer, entities.ValidationError(err))
		return
	}

	if err := json.Unmarshal(data, &document); err != nil {
		s.handleError(writer, entities.ValidationError(err))
		return
	}

	if err := document.Validate(); err != nil {
		s.handleError(writer, entities.ValidationError(err))
		return
	}

	id, cErr := actions.InsertDocument(ctx, s.repo, document)
	if cErr != nil {
		s.handleError(writer, cErr)
		return
	}

	if err := s.handleOK(writer, map[string]any{"id": id}); err != nil {
		s.handleError(writer, entities.ValidationError(err))
		return
	}
}
