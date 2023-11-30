package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
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
