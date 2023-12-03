package server

import (
	"net/http"

	"documents/internal/core"
	"documents/internal/library/web"
	"documents/internal/log"
	v1 "documents/internal/server/handlers/v1"
	"go.uber.org/zap"
)

type CommonHandler struct {
	Path     string
	Method   string
	Function v1.Handler
}

func (c CommonHandler) GetHandler(repo *core.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := c.Function(r, repo)
		if err != nil {
			log.Warn("error during request",
				zap.String("path", c.Path),
				zap.String("method", c.Method),
				zap.Any("error", err))

			if intErr := web.HandleError(w, err); intErr != nil {
				log.Warn("error handling request error",
					zap.String("path", c.Path),
					zap.String("method", c.Method),
					zap.Any("error", intErr))
			}

			return
		}

		if err := web.HandleOK(w, data); err != nil {
			log.Warn("error during handling response",
				zap.String("path", c.Path),
				zap.String("method", c.Method),
				zap.Any("error", err))

			if intErr := web.HandleError(w, err); intErr != nil {
				log.Warn("error handling request sending error",
					zap.String("path", c.Path),
					zap.String("method", c.Method),
					zap.Any("error", intErr))
			}

			return
		}

		log.Info("response sent",
			zap.String("path", c.Path),
			zap.String("method", c.Method),
			zap.Any("data", data))
	}
}

type AuthHandler struct {
	Path     string
	Method   string
	Function func(w http.ResponseWriter, r *http.Request, repo *core.Repository) error
}

func (c AuthHandler) GetHandler(repo *core.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := c.Function(w, r, repo)
		if err != nil {
			web.HandleError(w, err)
			return
		}
	}
}
