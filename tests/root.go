package tests

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"documents/internal/commands"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/golang-migrate/migrate/v4"
	"github.com/stretchr/testify/suite"
)

const (
	DefaultUser           = "testuser"
	DefaultPassword       = "password"
	DefaultDatabasePort   = 7654
	DefaultDatabase       = "documents"
	DefaultMigrationsPath = "file://../postgres/migrations"

	DefaultAppPort    = "7655"
	DefaultConfigPath = "../configs/dev.yaml"

	RequestTimeout = 5 * time.Second
)

var DefaultURL = fmt.Sprintf("http://localhost:%s", DefaultAppPort)
var DSN = fmt.Sprintf("postgresql://%s:%s@localhost:%d/%s?sslmode=disable",
	DefaultUser, DefaultPassword, DefaultDatabasePort, DefaultDatabase,
)

type BaseSuite struct {
	suite.Suite
	postgres *embeddedpostgres.EmbeddedPostgres
	command  *commands.Command
}

func (s *BaseSuite) migrate(dsn string) error {
	m, err := migrate.New(DefaultMigrationsPath, dsn)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		return err
	}

	return nil
}

func (s *BaseSuite) insertFakeUsers() {
	conn, err := sql.Open("postgres", DSN)
	s.Require().NoError(err)
	_, err = conn.Exec("INSERT INTO documents.t_user (username) VALUES ('spiderman')")
	s.Require().NoError(err)
	_, err = conn.Exec("INSERT INTO documents.t_user (username) VALUES ('ironman')")
	s.Require().NoError(err)
	_, err = conn.Exec("INSERT INTO documents.t_user (username) VALUES ('captain')")
	s.Require().NoError(err)
	_, err = conn.Exec("INSERT INTO documents.t_user (username) VALUES ('hawkeye')")
	s.Require().NoError(err)
	_, err = conn.Exec("INSERT INTO documents.t_user (username) VALUES ('black_widow')")
	s.Require().NoError(err)
	_, err = conn.Exec("INSERT INTO documents.t_user (username) VALUES ('antman')")
	s.Require().NoError(err)
	_, err = conn.Exec("INSERT INTO documents.t_user (username) VALUES ('loki')")
	s.Require().NoError(err)
	_, err = conn.Exec("INSERT INTO documents.t_user (username) VALUES ('tor')")
	s.Require().NoError(err)
	_, err = conn.Exec("INSERT INTO documents.t_user (username) VALUES ('captain_marvel')")
	s.Require().NoError(err)
	s.Require().NoError(conn.Close())
}

func (s *BaseSuite) authorize(userID int64) string {
	ctx := context.Background()
	ctx, err := s.command.Repository.SessionManager.Load(ctx, "")
	s.Require().NoError(err)
	s.command.Repository.SessionManager.Put(ctx, "user_id", userID)
	token, _, err := s.command.Repository.SessionManager.Commit(ctx)
	s.Require().NoError(err)

	return token
}

func (s *BaseSuite) SetupSuite() {
	var err error
	s.postgres = embeddedpostgres.NewDatabase(
		embeddedpostgres.DefaultConfig().
			Username(DefaultUser).
			Password(DefaultPassword).
			Database(DefaultDatabase).
			Port(DefaultDatabasePort),
	)
	err = s.postgres.Start()
	s.Require().NoError(err)

	defer func() {
		if err != nil {
			s.Require().NoError(s.postgres.Stop())
		}
	}()

	err = s.migrate(DSN)
	s.Require().NoError(err)

	s.insertFakeUsers()

	s.Require().NoError(os.Setenv("POSTGRES_DSN", DSN))
	s.Require().NoError(os.Setenv("PORT", DefaultAppPort))
	s.Require().NoError(os.Setenv("CONFIG", DefaultConfigPath))
	s.Require().NoError(os.Setenv("DBPASSWORD", DefaultPassword))

	var command commands.Command
	go func() {
		s.Require().NoError(command.Init())
		s.Require().ErrorIs(command.Start(), http.ErrServerClosed)
	}()
	s.command = &command

	time.Sleep(1 * time.Second)
}

func (s *BaseSuite) TearDownSuite() {
	s.Require().NoError(s.command.Cleanup())
	s.Require().NoError(s.postgres.Stop())
}

func sendRequest[OutputType any](
	ctx context.Context, method string, handler string,
	body map[string]any, query map[string]string, authCookie string,
) (respBody OutputType, resp *http.Response, err error) {
	ctx, cancel := context.WithTimeout(ctx, RequestTimeout)
	defer cancel()

	path, err := url.JoinPath(DefaultURL, handler)
	if err != nil {
		return respBody, nil, err
	}

	var req *http.Request

	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return respBody, nil, err
		}
		req, err = http.NewRequestWithContext(
			ctx, method, path, bytes.NewReader(bodyBytes),
		)
	} else {
		req, err = http.NewRequestWithContext(
			ctx, method, path, nil,
		)
	}
	if err != nil {
		return respBody, nil, err
	}

	q := req.URL.Query()
	for key, val := range query {
		q.Set(key, val)
	}
	req.URL.RawQuery = q.Encode()

	if authCookie != "" {
		req.AddCookie(&http.Cookie{
			Name:  "session",
			Value: authCookie,
		})
	}

	return doRequest[OutputType](req)
}

func doRequest[OutputType any](r *http.Request) (respBody OutputType, resp *http.Response, err error) {
	resp, err = http.DefaultClient.Do(r)
	if err != nil {
		return respBody, nil, err
	}

	defer func() {
		err = resp.Body.Close()
	}()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return respBody, nil, err
	}

	var result OutputType
	return result, resp, json.Unmarshal(b, &result)
}

func sendPost[OutputType any](
	ctx context.Context, handler string, body map[string]any, authCookie string,
) (respBody OutputType, resp *http.Response, err error) {
	return sendRequest[OutputType](ctx, http.MethodPost, handler, body, nil, authCookie)
}

func sendGet[OutputType any](
	ctx context.Context, handler string, query map[string]string, authCookie string,
) (respBody OutputType, resp *http.Response, err error) {
	return sendRequest[OutputType](ctx, http.MethodGet, handler, nil, query, authCookie)
}

func sendDelete[OutputType any](
	ctx context.Context, handler string, query map[string]string, authCookie string,
) (respBody OutputType, resp *http.Response, err error) {
	return sendRequest[OutputType](ctx, http.MethodDelete, handler, nil, query, authCookie)
}
