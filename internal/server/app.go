package server

import (
	"net/http"

	"github.com/novoseltcev/go-course/internal/model"
	"github.com/novoseltcev/go-course/internal/server/storage"
	"github.com/novoseltcev/go-course/internal/server/endpoints"
)

type Server struct {
	config Config
	counterStorage endpoints.MetricStorager[model.Counter]
	gaugeStorage endpoints.MetricStorager[model.Gauge]
}

func NewServer(config Config) *Server {
	return &Server{
		config: config,
		counterStorage: &storage.MemStorage[model.Counter]{Metrics: make(map[string]model.Counter)},
		gaugeStorage: &storage.MemStorage[model.Gauge]{Metrics: make(map[string]model.Gauge)},
	}
}

func (s *Server) Start() error {
	return http.ListenAndServe(s.config.Address, endpoints.GetRouter(&s.counterStorage, &s.gaugeStorage))
}
