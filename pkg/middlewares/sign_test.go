package middlewares_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/novoseltcev/go-course/pkg/middlewares"
	"github.com/novoseltcev/go-course/pkg/testutils"
	"github.com/novoseltcev/go-course/pkg/testutils/helpers"
)

func TestSignMiddleware(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(middlewares.Sign(testutils.STRING)(helpers.Webhook(t)))
	defer ts.Close()

	resp := helpers.SendRequest(t, ts, bytes.NewBufferString(testutils.JSON), nil)
	defer resp.Body.Close()
	body := helpers.ReadBody(t, resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.JSONEq(t, testutils.JSON, string(body))
	assert.NotEmptyf(t, resp.Header.Get("Hashsha256"), "Hashsha256 header not equal")
}
