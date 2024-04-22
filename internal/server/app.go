package server

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5"

	"github.com/jmoiron/sqlx"
	"github.com/novoseltcev/go-course/internal/model"
	"github.com/novoseltcev/go-course/internal/server/endpoints"
	"github.com/novoseltcev/go-course/internal/server/storage"
	
)


type Server struct {
	config Config
	db *sqlx.DB
	CounterStorage storage.MetricStorager[model.Counter]	`json:"counter"`
	GaugeStorage storage.MetricStorager[model.Gauge]		`json:"gauge"`
}


func NewServer(config Config) *Server {
	return &Server{
		config: config,
		db: nil,
		CounterStorage: nil,
		GaugeStorage: nil,
	}
}

func (s *Server) Start() error {
	db, err := sqlx.Open("pgx", s.config.DatabaseDsn)
	if err == nil {
		defer db.Close()
	}
	

	s.CounterStorage, s.GaugeStorage = storage.InitMetricStorages(db)

	if err := s.Restore(); err != nil {
		return err
	}
	
	if db == nil {
		go func() {
			for {
				time.Sleep(time.Duration(s.config.StoreInterval) * time.Second)
				s.Backup()
			}
		}()
	}
	

	if err := http.ListenAndServe(s.config.Address, endpoints.GetRouter(s.db, &s.CounterStorage, &s.GaugeStorage)); err != nil {
		return err
	}
	return s.Backup()
}

func (s *Server) Restore() error {
	if s.db != nil || !s.config.Restore || s.config.FileStoragePath == "" {
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
	if s.db != nil || s.config.FileStoragePath == "" {
		return nil
	}
	fd, err := os.OpenFile(s.config.FileStoragePath, os.O_WRONLY | os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer fd.Close()
	
	return json.NewEncoder(fd).Encode(s)
}
