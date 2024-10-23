package middlewares_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-course/pkg/middlewares"
)

const (
	testBody      = `{"ping": "pong"}`
	testSecretKey = "secret-key"
)

func testRequest(
	t *testing.T,
	ts *httptest.Server,
	method, path string,
	body io.Reader,
	headers map[string]string,
) (*http.Response, string) {
	t.Helper()

	req, err := http.NewRequestWithContext(context.Background(), method, ts.URL+path, body)
	require.NoError(t, err)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestSignMiddleware(t *testing.T) {
	t.Parallel()

	router := chi.NewRouter()
	router.Use(middlewares.Sign(testSecretKey))
	router.Post("/", http.HandlerFunc(testWebhook))
	ts := httptest.NewServer(router)

	defer ts.Close()

	resp, body := testRequest(t, ts, http.MethodPost, "/", bytes.NewBufferString(testBody), nil)

	defer resp.Body.Close()

	require.JSONEq(t, testBody, body)

	assert.NotEmptyf(t, resp.Header.Get("Hashsha256"), "Hashsha256 header not equal")
}
