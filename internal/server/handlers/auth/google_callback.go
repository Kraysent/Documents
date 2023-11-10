package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"documents/internal/actions"
	"documents/internal/actions/schema"
	"documents/internal/core"
	"documents/internal/library/web"
	"documents/internal/log"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const oauthGoogleUrlAPITemplate = "https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s"

func getGoogleConfig(redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  redirectURL,
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request, repo *core.Repository) (err error) {
	ctx := r.Context()

	oauthState, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthState.Value {
		return web.ValidationError(
			fmt.Errorf("invalid state in callback: got %s, expected %s", oauthState.Value, r.FormValue("state")),
		)
	}

	cfg := getGoogleConfig(repo.Config.Server.Callbacks.Google.RedirectURL)
	token, err := cfg.Exchange(ctx, r.FormValue("code"))
	if err != nil {
		return web.InternalError(err)
	}

	response, err := http.Get(fmt.Sprintf(oauthGoogleUrlAPITemplate, token.AccessToken))
	if err != nil {
		return web.InternalError(err)
	}
	defer func() {
		err = response.Body.Close()
	}()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return web.InternalError(err)
	}

	var request schema.GetOrCreateUserRequest

	if err := json.Unmarshal(contents, &request); err != nil {
		return web.InternalError(err)
	}

	if err := request.Validate(); err != nil {
		return web.ValidationError(err)
	}

	status, err := actions.GetOrCreateUser(ctx, repo, request)
	if err != nil {
		return err
	}

	log.Info("Obtained user info", zap.Any("status", status))

	http.Redirect(w, r, repo.Config.Server.Callbacks.BackRedirectURL, http.StatusTemporaryRedirect)

	return nil
}
