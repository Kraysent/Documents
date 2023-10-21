package server

import (
	"encoding/json"
	"net/http"

	"documents/internal/actions"
	"documents/internal/core"
	"documents/internal/log"
	"go.uber.org/zap"
)

type CommonHandler struct {
	Path     string
	Method   string
	Function func(*http.Request, *core.Repository) (any, error)
}

func (c CommonHandler) GetHandler(repo *core.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != c.Method {
			handle404(w)
			return
		}

		data, err := c.Function(r, repo)
		if err != nil {
			handleError(w, err)
			return
		}

		if err := handleOK(w, data); err != nil {
			handleError(w, err)
			return
		}
	}
}

func handle404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)

	data, intErr := json.Marshal(map[string]any{
		"message": "Not Found",
	})
	if intErr != nil {
		log.Warn("error during 404 response", zap.Error(intErr))
		return
	}

	if _, intErr := w.Write(data); intErr != nil {
		log.Warn("error during 404 response", zap.Error(intErr))
		return
	}
}

func handleError(writer http.ResponseWriter, err error) {
	codedErr, ok := err.(actions.CodedError)
	if !ok {
		codedErr = *actions.InternalError(err)
	}

	writer.WriteHeader(codedErr.HTTPCode)

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

func handleOK(writer http.ResponseWriter, data any) error {
	writer.WriteHeader(http.StatusOK)

	response, err := json.Marshal(NewResponse(data))
	if err != nil {
		return err
	}

	if _, err := writer.Write(response); err != nil {
		return err
	}

	log.Info("response sent", zap.Any("data", data))
	return nil
}
