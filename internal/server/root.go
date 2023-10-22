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

func (c CommonHandler) GetHandler(repo *core.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != c.Method {
			web.Handle404(w)
			return
		}

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
