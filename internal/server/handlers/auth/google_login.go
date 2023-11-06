package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"documents/internal/core"
	"documents/internal/library/web"
	"documents/internal/log"
	"go.uber.org/zap"
)

func GetGoogleLoginHandler(repo *core.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		expiration := time.Now().Add(365 * 24 * time.Hour)
		b := make([]byte, 16)

		_, err := rand.Read(b)
		if err != nil {
			web.HandleError(w, web.InternalError(err))
			return
		}

		state := base64.URLEncoding.EncodeToString(b)
		cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
		http.SetCookie(w, &cookie)

		url := getGoogleConfig(repo.Config.Server.Callbacks.Google.RedirectURL).AuthCodeURL(state)

		log.Info("Redirecting", zap.String("to", url))
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}
