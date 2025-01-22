package grpcserver_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/novoseltcev/go-course/pkg/grpcserver"
)

func TestNotify(t *testing.T) {
	t.Parallel()

	srv := grpcserver.New(":80")
	go srv.Run()

	require.NoError(t, srv.Shutdown(context.TODO()))
	assert.ErrorIs(t, <-srv.Notify(), grpc.ErrServerStopped)
}
