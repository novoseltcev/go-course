package helpers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// nolint: testifylint
func Webhook(t *testing.T) http.HandlerFunc {
	t.Helper()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")

		_, err = w.Write(body)
		require.NoError(t, err)
	})
}

func SendRequest(t *testing.T, ts *httptest.Server, body io.Reader, headers map[string]string) *http.Response {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, ts.URL, body)
	req.RequestURI = ""

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	return resp
}

func ReadBody(t *testing.T, body io.ReadCloser) []byte {
	t.Helper()

	b, err := io.ReadAll(body)
	require.NoError(t, err)

	return b
}
