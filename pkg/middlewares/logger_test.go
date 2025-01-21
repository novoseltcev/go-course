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

func TestLogger(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(middlewares.Logger(helpers.Webhook(t)))
	defer ts.Close()

	resp := helpers.SendRequest(t, ts, bytes.NewBufferString(testutils.JSON), nil)
	defer resp.Body.Close()
}

func TestLogger_WithWriteHeader(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(middlewares.Logger(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}),
	))
	defer ts.Close()

	resp := helpers.SendRequest(t, ts, bytes.NewBufferString(testutils.JSON), nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
