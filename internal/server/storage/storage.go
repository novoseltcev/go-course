package storage

import (
	"github.com/jmoiron/sqlx"

	"github.com/novoseltcev/go-course/internal/model"
	"github.com/novoseltcev/go-course/internal/server/storage/mem"
)


type MetricStorager[T model.Counter | model.Gauge] interface {
	GetByName(string) *T
	GetAll() []model.Metric[T]
	Update(string, T)
}


func InitMetricStorages(db *sqlx.DB) (MetricStorager[model.Counter], MetricStorager[model.Gauge]) {
	if db == nil {
		return &mem.Storage[model.Counter]{Metrics: make(map[string]model.Counter)}, &mem.Storage[model.Gauge]{Metrics: make(map[string]model.Gauge)}
	}

	return nil, nil
} 
