package tests

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type LinkSuite struct {
	BaseSuite
}

func TestLinkTestSuite(t *testing.T) {
	suite.Run(t, &LinkSuite{})
}

func (s *LinkSuite) TestCreateLinkExpiryBeforeNow() {
	token := s.authorize(1)
	ctx := context.Background()

	respBody, resp, err := sendPost[map[string]map[string]any](ctx, "/api/v1/document", map[string]any{
		"name":        "Some cool document name",
		"description": "even cooler description",
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	respBody, resp, err = sendPost[map[string]map[string]any](ctx, "/api/v1/link", map[string]any{
		"document_id": respBody["data"]["id"],
		"expiry_date": "2001-01-01T00:00:00Z",
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *LinkSuite) TestCreateLinkNoSuchDocument() {
	token := s.authorize(1)
	ctx := context.Background()

	_, resp, err := sendPost[map[string]map[string]any](ctx, "/api/v1/link", map[string]any{
		"document_id": "10101010-9fa0-9fa0-9fa0-9ce9ce9ce9ce",
		"expiry_date": "2500-01-01T00:00:00Z",
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *LinkSuite) TestCreateLinkWrongUser() {
	token := s.authorize(2)
	ctx := context.Background()

	respBody, resp, err := sendPost[map[string]map[string]any](ctx, "/api/v1/document", map[string]any{
		"name":        "Some cool document name",
		"description": "even cooler description",
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	token = s.authorize(3)
	_, resp, err = sendPost[map[string]map[string]any](ctx, "/api/v1/link", map[string]any{
		"document_id": respBody["data"]["id"],
		"expiry_date": "2500-01-01T00:00:00Z",
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *LinkSuite) TestCreateLinkHappyPath() {
	token := s.authorize(1)
	ctx := context.Background()

	respBody, resp, err := sendPost[map[string]map[string]any](ctx, "/api/v1/document", map[string]any{
		"name":        "Some cool document name",
		"description": "even cooler description",
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	respBody, resp, err = sendPost[map[string]map[string]any](ctx, "/api/v1/link", map[string]any{
		"document_id": respBody["data"]["id"],
		"expiry_date": "2500-01-01T00:00:00Z",
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	id, err := uuid.Parse(respBody["data"]["id"].(string))
	s.Require().NoError(err)
	s.Require().NotEqual(uuid.Nil, id)
}
