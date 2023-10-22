package auth

import (
	"context"
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

func getGoogleConfig(host string, port int) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  fmt.Sprintf("http://%s:%d/auth/google/callback", host, port),
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func GetGoogleCallbackHandler(repo *core.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()

		oauthState, _ := r.Cookie("oauthstate")

		if r.FormValue("state") != oauthState.Value {
			web.HandleError(w,
				web.ValidationError(
					fmt.Errorf("invalid state in callback: got %s, expected %s", oauthState.Value, r.FormValue("state"))),
			)
			return
		}

		token, err := getGoogleConfig(repo.Config.Server.Host, repo.Config.Server.Port).Exchange(ctx, r.FormValue("code"))
		if err != nil {
			web.HandleError(w, web.InternalError(err))
			return
		}

		response, err := http.Get(fmt.Sprintf(oauthGoogleUrlAPITemplate, token.AccessToken))
		if err != nil {
			web.HandleError(w, web.InternalError(err))
			return
		}
		defer response.Body.Close()

		contents, err := io.ReadAll(response.Body)
		if err != nil {
			web.HandleError(w, web.InternalError(err))
			return
		}

		var request schema.GetOrCreateUserRequest

		if err := json.Unmarshal(contents, &request); err != nil {
			web.HandleError(w, web.InternalError(err))
			return
		}

		if err := request.Validate(); err != nil {
			web.HandleError(w, web.ValidationError(err))
			return
		}

		status, err := actions.GetOrCreateUser(ctx, repo, request)
		if err != nil {
			web.HandleError(w, web.InternalError(err))
			return
		}

		log.Info("Obtained user info", zap.Any("status", status))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}
}
