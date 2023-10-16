package server

import (
	"encoding/json"
	"net/http"

	"documents/internal/core"
	"documents/internal/entities"
	"documents/internal/log"
	"go.uber.org/zap"
)

type Server struct {
	repo *core.Repository
}

func NewServer(repo *core.Repository) *Server {
	return &Server{repo: repo}
}

func (s *Server) handleError(writer http.ResponseWriter, err *entities.CodedError) {
	writer.WriteHeader(err.HTTPCode)

	data, intErr := json.Marshal(err)
	if intErr != nil {
		log.Warn(
			"error during handling another error",
			zap.Any("coded_error", err), zap.Error(intErr),
		)
		return
	}

	if _, intErr := writer.Write(data); intErr != nil {
		log.Warn(
			"error during handling another error",
			zap.Any("coded_error", err), zap.Error(intErr),
		)
		return
	}

	log.Warn("error during request", zap.Any("error", err))
}

func (s *Server) handleOK(writer http.ResponseWriter, data any) error {
	writer.WriteHeader(http.StatusOK)

	response, err := json.Marshal(entities.NewResponse(data))
	if err != nil {
		return err
	}

	if _, err := writer.Write(response); err != nil {
		return err
	}

	log.Info("response sent", zap.Any("data", data))
	return nil
}
