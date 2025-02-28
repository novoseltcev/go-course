package middlewares_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-course/pkg/middlewares"
	"github.com/novoseltcev/go-course/pkg/testutils"
	"github.com/novoseltcev/go-course/pkg/testutils/helpers"
)

func compress(t *testing.T, body []byte) []byte {
	t.Helper()

	buf := bytes.NewBuffer(nil)

	zb := gzip.NewWriter(buf)
	zb.Write(body)
	zb.Close()

	return buf.Bytes()
}

func uncompress(t *testing.T, body []byte) []byte {
	t.Helper()

	zr, err := gzip.NewReader(bytes.NewBuffer(body))
	require.NoError(t, err)

	b, err := io.ReadAll(zr)
	require.NoError(t, err)

	return b
}

func TestGzip_CompressedResponse_Success(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(middlewares.Gzip(helpers.Webhook(t)))
	defer ts.Close()

	resp := helpers.SendRequest(t, ts,
		bytes.NewBufferString(testutils.JSON),
		map[string]string{"Accept-Encoding": "gzip"},
	)
	defer resp.Body.Close()
	body := helpers.ReadBody(t, resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.JSONEq(t, testutils.JSON, string(uncompress(t, body)))
}

func TestGzip_CompressedRequest_Success(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(middlewares.Gzip(helpers.Webhook(t)))
	defer ts.Close()

	resp := helpers.SendRequest(t, ts,
		bytes.NewBuffer(compress(t, []byte(testutils.JSON))),
		map[string]string{"Content-Encoding": "gzip"},
	)
	defer resp.Body.Close()
	body := helpers.ReadBody(t, resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.JSONEq(t, testutils.JSON, string(body))
}

func TestGzip_WithoutCompression_Success(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "empty",
			value: "",
		},
		{
			name:  "unknown",
			value: "UNKNOWN",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ts := httptest.NewServer(middlewares.Gzip(helpers.Webhook(t)))
			defer ts.Close()

			resp := helpers.SendRequest(t, ts, bytes.NewBufferString(testutils.JSON), map[string]string{
				"Accept-Encoding":  tt.value,
				"Content-Encoding": tt.value,
			})
			defer resp.Body.Close()
			body := helpers.ReadBody(t, resp.Body)

			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.JSONEq(t, testutils.JSON, string(body))
		})
	}
}

func TestGzip_EmptyRequestWithContentEncoding_DecryptError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(middlewares.Gzip(helpers.Webhook(t)))
	defer ts.Close()

	resp := helpers.SendRequest(t, ts,
		nil,
		map[string]string{"Content-Encoding": "gzip"},
	)
	defer resp.Body.Close()
	body := helpers.ReadBody(t, resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, "cannot decompress body\n", string(body))
}
