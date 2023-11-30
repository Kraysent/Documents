package tests

import (
	"context"
	"database/sql"
	"net/http"
	"os"
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

func (s *DocumentSuite) insertFakeUsers() {
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

func (s *DocumentSuite) authorize(userID int64) string {
	ctx := context.Background()
	ctx, err := s.command.Repository.SessionManager.Load(ctx, "")
	s.Require().NoError(err)
	s.command.Repository.SessionManager.Put(ctx, "user_id", userID)
	token, _, err := s.command.Repository.SessionManager.Commit(ctx)
	s.Require().NoError(err)

	return token
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

	s.insertFakeUsers()

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
	s.Require().NoError(s.command.Cleanup())
	s.Require().NoError(s.postgres.Stop())
}

func (s *DocumentSuite) TestUnauthorizedCreation() {
	_, resp, err := sendPost[map[string]any](context.Background(), "/api/v1/document", map[string]any{
		"name":        "Some cool document name",
		"description": "even cooler description",
	}, "")
	s.Require().NoError(err)
	s.Require().Equal(http.StatusUnauthorized, resp.StatusCode)
}

func (s *DocumentSuite) TestInsertHappyPath() {
	token := s.authorize(1)
	ctx := context.Background()

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

func (s *DocumentSuite) TestInsertInvalidDocument() {
	token := s.authorize(1)
	ctx := context.Background()

	_, resp, err := sendPost[map[string]any](ctx, "/api/v1/document", map[string]any{
		"name":        "",
		"description": "the name is empty",
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *DocumentSuite) TestGetUserDocumentsUnauthorized() {
	ctx := context.Background()

	_, resp, err := sendGet[map[string]map[string]any](
		ctx, "/api/v1/user/documents", nil, "",
	)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusUnauthorized, resp.StatusCode)
}

func (s *DocumentSuite) TestGetUserDocumentsNoDocuments() {
	token := s.authorize(2)
	ctx := context.Background()

	respBody, resp, err := sendGet[map[string]map[string]any](
		ctx, "/api/v1/user/documents", nil, token,
	)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Require().Len(respBody["data"]["documents"].([]any), 0)
}

func (s *DocumentSuite) TestGetUserDocumentsHappyPath() {
	token := s.authorize(3)
	ctx := context.Background()

	_, resp, err := sendPost[map[string]any](ctx, "/api/v1/document", map[string]any{
		"name":        "Shield license",
		"description": "",
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	_, resp, err = sendPost[map[string]any](ctx, "/api/v1/document", map[string]any{
		"name":        "Old passport",
		"description": "",
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	respBody, resp, err := sendGet[map[string]map[string]any](
		ctx, "/api/v1/user/documents", nil, token,
	)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	docs := respBody["data"]["documents"].([]any)
	s.Require().Len(docs, 2)

	for _, doc := range docs {
		s.Require().NotEqual(uuid.Nil, doc.(map[string]any)["id"])
	}
}

func (s *DocumentSuite) TestGetDocumentByIDNoSuchDocument() {
	token := s.authorize(4)
	ctx := context.Background()

	_, resp, err := sendPost[map[string]any](ctx, "/api/v1/document", map[string]any{
		"name":        "Bow license",
		"description": "",
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	_, resp, err = sendGet[map[string]map[string]any](
		ctx, "/api/v1/document/id", map[string]string{"id": "10101010-9fa0-9fa0-9fa0-9ce9ce9ce9ce"}, token,
	)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *DocumentSuite) TestGetDocumentByIDWrongUser() {
	token := s.authorize(5)
	ctx := context.Background()

	respBody, resp, err := sendPost[map[string]map[string]any](ctx, "/api/v1/document", map[string]any{
		"name":        "Knife license",
		"description": "",
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	documentID := respBody["data"]["id"].(string)

	token = s.authorize(4)
	_, resp, err = sendGet[map[string]map[string]any](
		ctx, "/api/v1/document/id", map[string]string{"id": documentID}, token,
	)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *DocumentSuite) TestGetDocumentByIDHappyPath() {
	token := s.authorize(6)
	ctx := context.Background()

	respBody, resp, err := sendPost[map[string]map[string]any](ctx, "/api/v1/document", map[string]any{
		"name":        "Knife license",
		"description": "with a stamp!",
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	documentID := respBody["data"]["id"].(string)

	respBody, resp, err = sendGet[map[string]map[string]any](
		ctx, "/api/v1/document/id", map[string]string{"id": documentID}, token,
	)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Require().Equal("Knife license", respBody["data"]["name"])
}

func (s *DocumentSuite) TestDeleteDocumentThatDoesNotExist() {
	token := s.authorize(7)
	ctx := context.Background()

	_, resp, err := sendPost[map[string]any](ctx, "/api/v1/document", map[string]any{
		"name":        "Bow license",
		"description": "",
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	_, resp, err = sendDelete[map[string]map[string]any](
		ctx, "/api/v1/document", map[string]string{"id": "10101010-9fa0-9fa0-9fa0-9ce9ce9ce9ce"}, token,
	)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *DocumentSuite) TestDeleteDocumentWrongUser() {
	token := s.authorize(8)
	ctx := context.Background()

	respBody, resp, err := sendPost[map[string]map[string]any](ctx, "/api/v1/document", map[string]any{
		"name":        "Knife license",
		"description": "",
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	documentID := respBody["data"]["id"].(string)

	token = s.authorize(7)
	_, resp, err = sendDelete[map[string]map[string]any](
		ctx, "/api/v1/document", map[string]string{"id": documentID}, token,
	)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *DocumentSuite) TestDeleteDocumentHappyPath() {
	token := s.authorize(6)
	ctx := context.Background()

	respBody, resp, err := sendPost[map[string]map[string]any](ctx, "/api/v1/document", map[string]any{
		"name":        "Knife license",
		"description": "with a stamp!",
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	documentID := respBody["data"]["id"].(string)

	respBody, resp, err = sendDelete[map[string]map[string]any](
		ctx, "/api/v1/document", map[string]string{"id": documentID}, token,
	)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}
