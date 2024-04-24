package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/novoseltcev/go-course/internal/model"
	"github.com/novoseltcev/go-course/internal/server/endpoints"
	"github.com/novoseltcev/go-course/internal/server/storage"
	"github.com/novoseltcev/go-course/internal/server/storage/mem"
	"github.com/novoseltcev/go-course/internal/server/storage/pg"
)


type Server struct {
	config Config
	db *sqlx.DB
	MetricStorage storage.MetricStorager `json:"storage"`
}


func NewServer(config Config) *Server {
	return &Server{
		config: config,
		db: nil,
		MetricStorage: nil,
	}
}

func (s *Server) Start() error {
	if s.config.DatabaseDsn != "" {
		db, err := sqlx.Open("pgx", s.config.DatabaseDsn)
		if err != nil {
			return err
		}
		defer db.Close()
		s.db = db
		s.MetricStorage = &pg.Storage{DB: db}

		driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
		if err != nil {
			return err
		}
		m, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", driver)
		if err != nil {
			return err
		}
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			if downErr := m.Down(); downErr != nil {
				return errors.Join(err, downErr)
			}
			return err
		}
	} else {
		mapStorage := make(map[string]map[string]model.Metric)
		mapStorage["counter"] = make(map[string]model.Metric)
		mapStorage["gauge"] = make(map[string]model.Metric)
		s.MetricStorage = &mem.Storage{Metrics: mapStorage}
		
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

	return http.ListenAndServe(s.config.Address, endpoints.GetRouter(s.db, &s.MetricStorage))
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
