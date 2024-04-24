package storage

import (
	"context"

	"github.com/novoseltcev/go-course/internal/model"
)


type MetricStorager interface {
	GetByName(ctx context.Context, name, Type string) (*model.Metric, error)
	GetAll(ctx context.Context) ([]model.Metric, error)
	Save(ctx context.Context, metric model.Metric) error
	SaveAll(ctx context.Context, metrics []model.Metric) error
}
