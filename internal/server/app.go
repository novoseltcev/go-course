package server

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/novoseltcev/go-course/internal/server/endpoints"
	"github.com/novoseltcev/go-course/internal/server/middlewares"
	"github.com/novoseltcev/go-course/internal/server/storage"
	"github.com/novoseltcev/go-course/internal/server/storage/mem"
	"github.com/novoseltcev/go-course/internal/server/storage/pg"
	"net/http/pprof"
)


type Server struct {
	config Config
	MetricStorage storage.MetricStorager `json:"storage"`
}


func NewServer(config Config) *Server {
	return &Server{
		config: config,
		MetricStorage: nil,
	}
}

func (s *Server) Start() error {
	if s.config.DatabaseDsn != "" {
		storage, err := pg.New(s.config.DatabaseDsn)
		if err != nil {
			return nil
		}
		defer storage.Close()
		s.MetricStorage = storage
	} else {
		s.MetricStorage = mem.New()
		
		if err := s.Restore(); err != nil {
			return err
		}

		go func() {
			for {
				time.Sleep(time.Duration(s.config.StoreInterval) * time.Second)
				s.Backup()
			}
		}()

		defer s.Backup()
	}

	return http.ListenAndServe(s.config.Address, s.GetRouter())
}

func (s *Server) Restore() error {
	if !s.config.Restore || s.config.FileStoragePath == "" {
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


func (s *Server) GetRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middlewares.Logger)

	if s.config.SecretKey != "" {
		r.Use(middlewares.CheckSum(s.config.SecretKey))
	}
	r.Use(middlewares.Gzip)
	if s.config.SecretKey != "" {
		r.Use(middlewares.Sign(s.config.SecretKey))
	}

	storage := s.MetricStorage

	r.Get(`/ping`, endpoints.Ping(storage))
	r.Get(`/`, endpoints.Index(storage))
	r.Post(`/update/{metricType}/{metricName}/{metricValue}`, endpoints.UpdateMetric(storage))
	r.Get(`/value/{metricType}/{metricName}`, endpoints.GetOneMetric(storage))
	r.Post(`/update/`, endpoints.UpdateMetricFromJSON(storage))
	r.Post(`/value/`, endpoints.GetOneMetricFromJSON(storage))
	r.Post(`/updates/`, endpoints.UpdateMetricsBatch(storage))
	
    r.HandleFunc("/debug/pprof", pprof.Index)
    r.HandleFunc("/debug/pprof/profile", pprof.Profile)
    r.HandleFunc("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)

	return r
}
