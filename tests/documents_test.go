package tests

import (
	"context"
	"net/http"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type DocumentSuite struct {
	BaseSuite
}

func TestDocumentTestSuite(t *testing.T) {
	suite.Run(t, &DocumentSuite{})
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
