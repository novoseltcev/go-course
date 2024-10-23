package middlewares_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-course/pkg/middlewares"
)

func TestChecksumMiddlewareSuccess(t *testing.T) {
	t.Parallel()

	router := chi.NewRouter()
	router.Use(middlewares.CheckSum(testSecretKey))
	router.Post("/", http.HandlerFunc(testWebhook))
	ts := httptest.NewServer(router)

	defer ts.Close()

	resp, _ := testRequest(
		t,
		ts,
		http.MethodPost,
		"/",
		bytes.NewBufferString(testBody),
		map[string]string{
			// computed by site: https://www.devglan.com/online-tools/hmac-sha256-online
			"Hashsha256": "4a7991274e9f3faa4f7fa5828b25ac896532438fdeb63e34c38431b75cbc03e4",
		},
	)

	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestChecksumMiddlewareErr(t *testing.T) {
	t.Parallel()

	router := chi.NewRouter()
	router.Use(middlewares.CheckSum(testSecretKey))
	router.Post("/", http.HandlerFunc(testWebhook))
	ts := httptest.NewServer(router)

	defer ts.Close()

	resp, _ := testRequest(
		t,
		ts,
		http.MethodPost,
		"/",
		bytes.NewBufferString(testBody),
		map[string]string{
			"Hashsha256": "invalid-hash",
		},
	)

	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
