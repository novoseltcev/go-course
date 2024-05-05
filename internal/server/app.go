package server

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/novoseltcev/go-course/internal/server/endpoints"
	"github.com/novoseltcev/go-course/internal/server/storage"
	"github.com/novoseltcev/go-course/internal/server/storage/mem"
	"github.com/novoseltcev/go-course/internal/server/storage/pg"
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

	return http.ListenAndServe(s.config.Address, endpoints.GetRouter(s.MetricStorage))
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
