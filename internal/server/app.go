package server

import (
	"net/http"

	"github.com/novoseltcev/go-course/internal/server/endpoints"
	"github.com/novoseltcev/go-course/internal/server/storage"
	"github.com/novoseltcev/go-course/internal/types"
)

type Server struct {
	config Config
	counterStorage endpoints.MetricStorager[types.Counter]
	gaugeStorage endpoints.MetricStorager[types.Gauge]
}

func NewServer(config Config) *Server {
	return &Server{
		config: config,
		counterStorage: &storage.MemStorage[types.Counter]{Metrics: make(map[string]types.Counter)},
		gaugeStorage: &storage.MemStorage[types.Gauge]{Metrics: make(map[string]types.Gauge)},
	}
}

func (s *Server) Start() error {
	return http.ListenAndServe(s.config.Address, endpoints.GetRouter(&s.counterStorage, &s.gaugeStorage))
}
