package web

import (
	"encoding/json"
	"net/http"

	"documents/internal/log"
	"go.uber.org/zap"
)

func Handle404(w http.ResponseWriter) {
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

func HandleError(writer http.ResponseWriter, err error) error {
	codedErr, ok := err.(CodedError)
	if !ok {
		codedErr = *InternalError(err)
	}

	writer.WriteHeader(codedErr.HTTPCode)

	data, intErr := json.Marshal(err)
	if intErr != nil {
		return intErr
	}

	if _, intErr := writer.Write(data); intErr != nil {
		return intErr
	}

	return nil
}

func HandleOK(writer http.ResponseWriter, data any) error {
	writer.WriteHeader(http.StatusOK)

	response, err := json.Marshal(NewResponse(data))
	if err != nil {
		return err
	}

	if _, err := writer.Write(response); err != nil {
		return err
	}

	return nil
}
