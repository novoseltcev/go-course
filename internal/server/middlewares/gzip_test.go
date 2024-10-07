package middlewares_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-course/internal/server/middlewares"
)

func webhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	defer r.Body.Close()

	w.Write(body)
	w.Header().Set("Content-Type", "application/json")
}

const testGzipCompression = `{"ping": "pong"}`

func TestGzipCompressionWithoutAcceptEncoding(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(middlewares.Gzip(http.HandlerFunc(webhook)))

	defer srv.Close()

	buf := bytes.NewBufferString(testGzipCompression)
	r := httptest.NewRequest(http.MethodPost, srv.URL, buf)
	r.RequestURI = ""
	r.Header.Set("Accept-Encoding", "")

	resp, err := http.DefaultClient.Do(r)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.JSONEq(t, testGzipCompression, string(b))
}

func TestGzipCompressionSend(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(middlewares.Gzip(http.HandlerFunc(webhook)))

	defer srv.Close()

	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	_, err := zb.Write([]byte(testGzipCompression))
	require.NoError(t, err)
	err = zb.Close()

	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, srv.URL, buf)
	r.RequestURI = ""
	r.Header.Set("Content-Encoding", "gzip")
	r.Header.Set("Accept-Encoding", "")

	resp, err := http.DefaultClient.Do(r)
	require.NoError(t, err)

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.JSONEq(t, testGzipCompression, string(b))
}

func TestGzipCompressionAccept(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(middlewares.Gzip(http.HandlerFunc(webhook)))

	defer srv.Close()

	buf := bytes.NewBufferString(testGzipCompression)
	r := httptest.NewRequest(http.MethodPost, srv.URL, buf)
	r.RequestURI = ""

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(r)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()

	zr, err := gzip.NewReader(resp.Body)
	require.NoError(t, err)

	b, err := io.ReadAll(zr)
	require.NoError(t, err)
	require.JSONEq(t, testGzipCompression, string(b))
}
