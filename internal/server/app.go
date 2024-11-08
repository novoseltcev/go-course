package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/pprof"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/server/endpoints"
	"github.com/novoseltcev/go-course/internal/storages"
	"github.com/novoseltcev/go-course/pkg/middlewares"
)

type Server struct {
	config        *Config
	MetricStorage storages.MetricStorager `json:"storage"`
}

func NewServer(config *Config) *Server {
	return &Server{
		config:        config,
		MetricStorage: nil,
	}
}

func (s *Server) Start() error {
	if s.config.DatabaseDsn == "" {
		s.MetricStorage = storages.NewMemStorage()
	} else {
		storage := storages.NewPgStorage(s.config.DatabaseDsn)
		defer storage.Close()

		s.MetricStorage = storage
	}

	if s.config.DatabaseDsn == "" {
		if err := s.Restore(); err != nil {
			return err
		}

		defer s.Backup() //nolint:errcheck
	}

	if s.config.DatabaseDsn == "" && s.config.StoreInterval > 0 {
		ctx, cancel := context.WithCancel(context.Background())

		defer cancel()

		go s.BackupWorker(ctx)
	}

	return http.ListenAndServe(s.config.Address, s.GetRouter())
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

	r.Mount("/", endpoints.NewAPIRouter(s.MetricStorage))

	r.HandleFunc("/debug/pprof", pprof.Index)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)

	return r
}

func (s *Server) Restore() error {
	if !s.config.Restore || s.config.FileStoragePath == "" {
		return nil
	}

	fd, err := os.OpenFile(s.config.FileStoragePath, os.O_RDONLY, 0o666)
	if os.IsNotExist(err) {
		_, err := os.OpenFile(s.config.FileStoragePath, os.O_CREATE, 0o666)

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

	fd, err := os.OpenFile(s.config.FileStoragePath, os.O_WRONLY|os.O_CREATE, 0o666)
	if err != nil {
		return err
	}

	defer fd.Close()

	if err := json.NewEncoder(fd).Encode(s); err != nil {
		log.WithError(err).Error("failed to backup metrics")

		return err
	}

	return nil
}

func (s *Server) BackupWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Duration(s.config.StoreInterval) * time.Second)
		}

		if err := s.Backup(); err != nil {
			log.WithError(err).Error("failed to backup metrics")
		}
	}
}
