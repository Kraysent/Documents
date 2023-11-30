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
	"syscall"
	"testing"
	"time"

	"documents/internal/commands"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
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

	RequestTimeout = 1 * time.Second
)

var DefaultURL = fmt.Sprintf("http://localhost:%s", DefaultAppPort)
var DSN = fmt.Sprintf("postgresql://%s:%s@localhost:%d/%s?sslmode=disable",
	DefaultUser, DefaultPassword, DefaultDatabasePort, DefaultDatabase,
)

type DocumentSuite struct {
	suite.Suite
	postgres *embeddedpostgres.EmbeddedPostgres
	command  *commands.Command
}

func TestDocumentTestSuite(t *testing.T) {
	suite.Run(t, &DocumentSuite{})
}

func (s *DocumentSuite) migrate(dsn string) error {
	m, err := migrate.New(DefaultMigrationsPath, dsn)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		return err
	}

	return nil
}

func sendPost[T any](
	ctx context.Context, handler string, data map[string]any, authCookie string,
) (respBody T, resp *http.Response, err error) {
	ctx, cancel := context.WithTimeout(ctx, RequestTimeout)
	defer cancel()

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return respBody, nil, err
	}

	path, err := url.JoinPath(DefaultURL, handler)
	if err != nil {
		return respBody, nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx, http.MethodPost, path, bytes.NewReader(dataBytes),
	)
	if err != nil {
		return respBody, nil, err
	}

	if authCookie != "" {
		req.AddCookie(&http.Cookie{
			Name:  "session",
			Value: authCookie,
		})
	}

	resp, err = http.DefaultClient.Do(req)
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

	var result T
	return result, resp, json.Unmarshal(b, &result)
}

func (s *DocumentSuite) SetupSuite() {
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

	s.Require().NoError(os.Setenv("POSTGRES_DSN", DSN))
	s.Require().NoError(os.Setenv("PORT", DefaultAppPort))
	s.Require().NoError(os.Setenv("CONFIG", DefaultConfigPath))
	s.Require().NoError(os.Setenv("DBPASSWORD", DefaultPassword))

	var command commands.Command
	s.Require().NoError(command.Init())
	go func() {
		s.Require().NoError(command.Start())
	}()
	s.command = &command

	time.Sleep(1 * time.Second)
}

func (s *DocumentSuite) TearDownSuite() {
	s.Require().NoError(s.postgres.Stop())
	p, _ := os.FindProcess(syscall.Getpid())
	s.Require().NoError(p.Signal(syscall.SIGINT))
}

func (s *DocumentSuite) TestUnauthorizedCreation() {
	_, resp, err := sendPost[map[string]any](context.Background(), "/api/v1/document", map[string]any{
		"name":        "Some cool document name",
		"description": "even cooler description",
	}, "")
	s.Require().NoError(err)
	s.Require().Equal(http.StatusUnauthorized, resp.StatusCode)
}

func (s *DocumentSuite) TestCreation() {
	conn, err := sql.Open("postgres", DSN)
	s.Require().NoError(err)
	_, err = conn.Exec("INSERT INTO documents.t_user (username) VALUES ('spiderman')")
	s.Require().NoError(err)
	_, err = conn.Exec("INSERT INTO documents.t_user (username) VALUES ('ironman')")
	s.Require().NoError(err)
	_, err = conn.Exec("INSERT INTO documents.t_user (username) VALUES ('captain')")
	s.Require().NoError(err)

	ctx := context.Background()
	ctx, err = s.command.Repository.SessionManager.Load(ctx, "")
	s.Require().NoError(err)
	s.command.Repository.SessionManager.Put(ctx, "user_id", int64(1))
	token, _, err := s.command.Repository.SessionManager.Commit(ctx)
	s.Require().NoError(err)

	respBody, resp, err := sendPost[map[string]map[string]any](ctx, "/api/v1/document", map[string]any{
		"name":        "Some cool document name",
		"description": "even cooler description",
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	id, err := uuid.Parse(respBody["data"]["id"].(string))
	s.Require().NoError(err)
	s.Require().NotEqual(uuid.Nil, id)
}
