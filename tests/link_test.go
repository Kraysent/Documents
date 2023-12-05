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

func (s *LinkSuite) TestDeleteDocumentWithLink() {
	token := s.authorize(3)
	ctx := context.Background()

	respBody, resp, err := sendPost[map[string]map[string]any](ctx, "/api/v1/document", map[string]any{
		"name":        "Some cool document name",
		"description": "even cooler description",
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	documentID := respBody["data"]["id"].(string)

	respBody, resp, err = sendPost[map[string]map[string]any](ctx, "/api/v1/link", map[string]any{
		"document_id": documentID,
		"expiry_date": time.Now().Add(10 * time.Hour).Format(time.RFC3339),
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	respBody, resp, err = sendDelete[map[string]map[string]any](
		ctx, "/api/v1/document", map[string]string{"id": documentID}, token,
	)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *LinkSuite) addDefaultLinks(ctx context.Context, token string, enabledNumber, disabledNumber int) (documentID string) {
	respBody, resp, err := sendPost[map[string]map[string]any](ctx, "/api/v1/document", map[string]any{
		"name":        "Some cool document name",
		"description": "even cooler description",
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	documentID = respBody["data"]["id"].(string)

	for i := 0; i < enabledNumber; i++ {
		respBody, resp, err = sendPost[map[string]map[string]any](ctx, "/api/v1/link", map[string]any{
			"document_id": documentID,
			"expiry_date": time.Now().Add(time.Duration(enabledNumber-i) * time.Hour).Format(time.RFC3339),
		}, token)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)
	}

	for i := 0; i < disabledNumber; i++ {
		respBody, resp, err = sendPost[map[string]map[string]any](ctx, "/api/v1/link", map[string]any{
			"document_id": documentID,
			"expiry_date": time.Now().Add(time.Duration(disabledNumber-i) * time.Hour).Format(time.RFC3339),
		}, token)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)

		_, resp, err = sendDelete[map[string]map[string]any](ctx, "/api/v1/link", map[string]string{
			"id": respBody["data"]["id"].(string),
		}, token)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)
	}

	return documentID
}

func (s *LinkSuite) TestGetLinksPageSizeTooBig() {
	token := s.authorize(4)
	ctx := context.Background()

	documentID := s.addDefaultLinks(ctx, token, 10, 0)

	_, resp, err := sendGet[map[string]map[string]any](ctx, "/api/v1/links", map[string]string{
		"document_id": documentID,
		"page_size":   "1000",
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *LinkSuite) TestGetLinksInvalidStatus() {
	token := s.authorize(5)
	ctx := context.Background()

	documentID := s.addDefaultLinks(ctx, token, 10, 0)

	_, resp, err := sendGet[map[string]map[string]any](ctx, "/api/v1/links", map[string]string{
		"document_id": documentID,
		"page_size":   "10",
		"status":      "very_much_nonexistent_status",
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *LinkSuite) TestGetLinksWrongUser() {
	token := s.authorize(6)
	ctx := context.Background()

	documentID := s.addDefaultLinks(ctx, token, 10, 0)

	token = s.authorize(5)
	_, resp, err := sendGet[map[string]map[string]any](ctx, "/api/v1/links", map[string]string{
		"document_id": documentID,
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *LinkSuite) TestGetLinksHappyPathAllResultsOnSinglePage() {
	token := s.authorize(7)
	ctx := context.Background()

	documentID := s.addDefaultLinks(ctx, token, 10, 0)
	respBody, resp, err := sendGet[map[string]map[string]any](ctx, "/api/v1/links", map[string]string{
		"document_id": documentID,
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	res := respBody["data"]["links"].([]any)
	s.Require().Len(res, 10)
	for _, linkAny := range res {
		link := linkAny.(map[string]any)
		id := link["id"].(string)
		_, err := uuid.Parse(id)
		s.Require().NoError(err)
	}
}

func (s *LinkSuite) TestGetLinksHappyPathAllResultsOnMultiplePages() {
	token := s.authorize(8)
	ctx := context.Background()

	documentID := s.addDefaultLinks(ctx, token, 30, 0)
	respBody, resp, err := sendGet[map[string]map[string]any](ctx, "/api/v1/links", map[string]string{
		"document_id": documentID,
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	res := respBody["data"]["links"].([]any)
	s.Require().Len(res, 25)
	for _, linkAny := range res {
		link := linkAny.(map[string]any)
		id := link["id"].(string)
		_, err := uuid.Parse(id)
		s.Require().NoError(err)
	}

	respBody, resp, err = sendGet[map[string]map[string]any](ctx, "/api/v1/links", map[string]string{
		"document_id": documentID,
		"page":        "1",
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	res = respBody["data"]["links"].([]any)
	s.Require().Len(res, 5)
	for _, linkAny := range res {
		link := linkAny.(map[string]any)
		id := link["id"].(string)
		_, err := uuid.Parse(id)
		s.Require().NoError(err)
	}
}

func (s *LinkSuite) TestGetLinksHappyPathOnlyDisabledLinks() {
	token := s.authorize(9)
	ctx := context.Background()

	documentID := s.addDefaultLinks(ctx, token, 10, 8)
	respBody, resp, err := sendGet[map[string]map[string]any](ctx, "/api/v1/links", map[string]string{
		"document_id": documentID,
		"status":      "disabled",
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	res := respBody["data"]["links"].([]any)
	s.Require().Len(res, 8)
	for _, linkAny := range res {
		link := linkAny.(map[string]any)
		id := link["id"].(string)
		_, err := uuid.Parse(id)
		s.Require().NoError(err)
	}
}

func (s *LinkSuite) TestGetLinksHappyPathNoLinks() {
	token := s.authorize(9)
	ctx := context.Background()

	documentID := s.addDefaultLinks(ctx, token, 0, 0)
	respBody, resp, err := sendGet[map[string]map[string]any](ctx, "/api/v1/links", map[string]string{
		"document_id": documentID,
	}, token)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	res := respBody["data"]["links"].([]any)
	s.Require().Len(res, 0)
}
