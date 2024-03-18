package server

import (
	"net/http"
	"github.com/novoseltcev/go-course/internal/server/endpoints"
	"github.com/novoseltcev/go-course/internal/server/storage"
)

type Server struct {
	config Config
	counterStorage storage.Storage[storage.Counter]
	gaugeStorage storage.Storage[storage.Gauge]
}

func NewServer(config Config) *Server {
	return &Server{
		config: config,
		counterStorage: storage.MemStorage[storage.Counter]{Metrics: make(map[string]storage.Counter)},
		gaugeStorage: storage.MemStorage[storage.Gauge]{Metrics: make(map[string]storage.Gauge)},
	}
}

func (s *Server) Start() error {
	return http.ListenAndServe(s.config.Address, endpoints.GetRouter(&s.counterStorage, &s.gaugeStorage))
}
