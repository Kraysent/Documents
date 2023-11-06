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

func GoogleLoginHandler(w http.ResponseWriter, r *http.Request, repo *core.Repository) error {
	expiration := time.Now().Add(365 * 24 * time.Hour)
	b := make([]byte, 16)

	_, err := rand.Read(b)
	if err != nil {
		return web.InternalError(err)
	}

	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	url := getGoogleConfig(repo.Config.Server.Callbacks.Google.RedirectURL).AuthCodeURL(state)

	log.Info("Redirecting", zap.String("to", url))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

	return nil
}
