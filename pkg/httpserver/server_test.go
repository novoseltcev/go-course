package httpserver_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/novoseltcev/go-course/pkg/httpserver"
)

func TestNew(t *testing.T) {
	t.Parallel()

	httpserver.New(nil)
}

func TestNewWithAllOpts(t *testing.T) {
	t.Parallel()

	httpserver.New(nil,
		httpserver.WithAddr(""),
		httpserver.WithShutdownTimeout(0),
		httpserver.WithReadTimeout(0),
		httpserver.WithWriteTimeout(0),
	)
}

func TestNotify(t *testing.T) {
	t.Parallel()

	srv := httpserver.New(nil)

	require.NoError(t, srv.Shutdown())
	assert.ErrorIs(t, <-srv.Notify(), http.ErrServerClosed)
}
