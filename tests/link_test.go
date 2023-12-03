package tests

import (
	"context"
	"net/http"
	"testing"
	"time"

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
		"expiry_date": time.Now().Add(-5 * time.Hour).Format(time.RFC3339),
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *LinkSuite) TestCreateLinkNoSuchDocument() {
	token := s.authorize(1)
	ctx := context.Background()

	_, resp, err := sendPost[map[string]map[string]any](ctx, "/api/v1/link", map[string]any{
		"document_id": "10101010-9fa0-9fa0-9fa0-9ce9ce9ce9ce",
		"expiry_date": time.Now().Add(5 * time.Hour).Format(time.RFC3339),
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
		"expiry_date": time.Now().Add(5 * time.Hour).Format(time.RFC3339),
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

func (s *LinkSuite) TestGetDocumentByLinkExpiryPassed() {
	token := s.authorize(3)
	ctx := context.Background()

	respBody, resp, err := sendPost[map[string]map[string]any](ctx, "/api/v1/document", map[string]any{
		"name":        "Some cool document name",
		"description": "even cooler description",
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	respBody, resp, err = sendPost[map[string]map[string]any](ctx, "/api/v1/link", map[string]any{
		"document_id": respBody["data"]["id"],
		"expiry_date": time.Now().Add(1 * time.Second).Format(time.RFC3339),
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	time.Sleep(2 * time.Second)

	_, resp, err = sendGet[map[string]map[string]any](ctx, "/api/v1/link", map[string]string{
		"id": respBody["data"]["id"].(string),
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *LinkSuite) TestGetDocumentByLinkHappyPath() {
	token := s.authorize(3)
	ctx := context.Background()

	respBody, resp, err := sendPost[map[string]map[string]any](ctx, "/api/v1/document", map[string]any{
		"name":        "Some cool document name",
		"description": "even cooler description",
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	respBody, resp, err = sendPost[map[string]map[string]any](ctx, "/api/v1/link", map[string]any{
		"document_id": respBody["data"]["id"],
		"expiry_date": time.Now().Add(10 * time.Hour).Format(time.RFC3339),
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	respBody, resp, err = sendGet[map[string]map[string]any](ctx, "/api/v1/link", map[string]string{
		"id": respBody["data"]["id"].(string),
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Require().Equal("Some cool document name", respBody["data"]["name"])
}

func (s *LinkSuite) TestGetDocumentByLinkDisabledLink() {
	token := s.authorize(3)
	ctx := context.Background()

	respBody, resp, err := sendPost[map[string]map[string]any](ctx, "/api/v1/document", map[string]any{
		"name":        "Some cool document name",
		"description": "even cooler description",
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	respBody, resp, err = sendPost[map[string]map[string]any](ctx, "/api/v1/link", map[string]any{
		"document_id": respBody["data"]["id"],
		"expiry_date": time.Now().Add(10 * time.Hour).Format(time.RFC3339),
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	_, resp, err = sendDelete[map[string]map[string]any](ctx, "/api/v1/link", map[string]string{
		"id": respBody["data"]["id"].(string),
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	respBody, resp, err = sendGet[map[string]map[string]any](ctx, "/api/v1/link", map[string]string{
		"id": respBody["data"]["id"].(string),
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *LinkSuite) TestDisableLinkNoSuchLink() {
	token := s.authorize(3)
	ctx := context.Background()

	_, resp, err := sendDelete[map[string]map[string]any](ctx, "/api/v1/link", map[string]string{
		"id": "10101010-9fa0-9fa0-9fa0-9ce9ce9ce9ce",
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *LinkSuite) TestDisableAlreadyDisabledLink() {
	token := s.authorize(3)
	ctx := context.Background()

	respBody, resp, err := sendPost[map[string]map[string]any](ctx, "/api/v1/document", map[string]any{
		"name":        "Some cool document name",
		"description": "even cooler description",
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	respBody, resp, err = sendPost[map[string]map[string]any](ctx, "/api/v1/link", map[string]any{
		"document_id": respBody["data"]["id"],
		"expiry_date": time.Now().Add(10 * time.Hour).Format(time.RFC3339),
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	_, resp, err = sendDelete[map[string]map[string]any](ctx, "/api/v1/link", map[string]string{
		"id": respBody["data"]["id"].(string),
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	_, resp, err = sendDelete[map[string]map[string]any](ctx, "/api/v1/link", map[string]string{
		"id": respBody["data"]["id"].(string),
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}
