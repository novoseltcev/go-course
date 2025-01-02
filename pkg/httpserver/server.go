// Package httpserver implements HTTP server.
package httpserver

import (
	"context"
	"net/http"
	"time"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultAddr            = ":80"
	defaultShutdownTimeout = 3 * time.Second
)

// Server represents HTTP server.
type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

// New creates new server and start listening.
func New(handler http.Handler, opts ...Option) *Server {
	srv := &Server{
		server: &http.Server{
			Handler:      handler,
			ReadTimeout:  defaultReadTimeout,
			WriteTimeout: defaultWriteTimeout,
			Addr:         defaultAddr,
		},
		notify:          make(chan error, 1),
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(srv)
	}

	go func() {
		srv.notify <- srv.server.ListenAndServe()
		close(srv.notify)
	}()

	return srv
}

// Notify returns a channel that will be closed when server is stopped.
func (srv *Server) Notify() <-chan error {
	return srv.notify
}

// Shutdown server gracefully with timeout.
func (srv *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), srv.shutdownTimeout)
	defer cancel()

	return srv.server.Shutdown(ctx)
}
