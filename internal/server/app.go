package server

import (
	"os"
	"time"
	"net/http"
	"encoding/json"

	"github.com/novoseltcev/go-course/internal/model"
	"github.com/novoseltcev/go-course/internal/server/endpoints"
	"github.com/novoseltcev/go-course/internal/server/storage"
)


type Server struct {
	config Config
	CounterStorage endpoints.MetricStorager[model.Counter]	`json:"counter"`
	GaugeStorage endpoints.MetricStorager[model.Gauge]		`json:"gauge"`
}


func NewServer(config Config) *Server {
	return &Server{
		config: config,
		CounterStorage: &storage.MemStorage[model.Counter]{Metrics: make(map[string]model.Counter)},
		GaugeStorage: &storage.MemStorage[model.Gauge]{Metrics: make(map[string]model.Gauge)},
	}
}

func (s *Server) Start() error {
	if s.config.Restore {
		if err := s.Restore(); err != nil {
			return err
		}
	}

	go func() {
		for {
			time.Sleep(time.Duration(s.config.StoreInterval) * time.Second)
			s.Backup()
		}
	}()

	if err := http.ListenAndServe(s.config.Address, endpoints.GetRouter(&s.CounterStorage, &s.GaugeStorage)); err != nil {
		return err
	}
	return s.Backup()
}

func (s *Server) Restore() error {
	if s.config.FileStoragePath == "" {
		return nil
	}

	fd, err := os.OpenFile(s.config.FileStoragePath, os.O_RDONLY, 0666)
	if os.IsNotExist(err) {
		_, err := os.OpenFile(s.config.FileStoragePath, os.O_CREATE, 0666)
		return err
	}
	if err != nil {
		return err
	}
	defer fd.Close()

	return json.NewDecoder(fd).Decode(s)
}

func (s *Server) Backup() error {
	if s.config.FileStoragePath == "" {
		return nil
	}
	fd, err := os.OpenFile(s.config.FileStoragePath, os.O_WRONLY | os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer fd.Close()
	
	return json.NewEncoder(fd).Encode(s)
}
