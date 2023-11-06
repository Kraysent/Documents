package server

import (
	"net/http"

	"documents/internal/core"
	"documents/internal/library/web"
)

type CommonHandler struct {
	Path     string
	Method   string
	Function func(*http.Request, *core.Repository) (any, error)
}

func (c CommonHandler) GetHandler(repo *core.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := c.Function(r, repo)
		if err != nil {
			web.HandleError(w, err)
			return
		}

		if err := web.HandleOK(w, data); err != nil {
			web.HandleError(w, err)
			return
		}
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
