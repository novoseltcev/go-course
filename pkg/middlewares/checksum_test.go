package middlewares_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-course/pkg/middlewares"
	"github.com/novoseltcev/go-course/pkg/testutils"
	"github.com/novoseltcev/go-course/pkg/testutils/helpers"
)

func TestChecksumMiddlewareSuccess(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(middlewares.CheckSum("secret")(helpers.Webhook(t)))
	defer ts.Close()

	resp := helpers.SendRequest(t, ts, bytes.NewBufferString(testutils.JSON), map[string]string{
		// computed by site: https://www.devglan.com/online-tools/hmac-sha256-online
		"Hashsha256": "23fff24c4b9835c6179de19103c6c640150d07d8a72c987b030b541a9d988736",
	})
	defer resp.Body.Close()
	body := helpers.ReadBody(t, resp.Body)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, testutils.JSON, string(body))
}

func TestChecksumMiddlewareErr(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(middlewares.CheckSum(testutils.STRING)(helpers.Webhook(t)))
	defer ts.Close()

	resp := helpers.SendRequest(t, ts,
		bytes.NewBufferString(testutils.JSON),
		map[string]string{"Hashsha256": "invalid"},
	)
	defer resp.Body.Close()
	body := helpers.ReadBody(t, resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, string(body), "encoding/hex: invalid byte")
}
