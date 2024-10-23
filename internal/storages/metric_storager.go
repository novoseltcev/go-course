package storages

import (
	"context"

	"github.com/novoseltcev/go-course/internal/schemas"
)

type MetricStorager interface {
	GetByName(ctx context.Context, name, Type string) (*schemas.Metric, error)
	GetAll(ctx context.Context) ([]schemas.Metric, error)
	Save(ctx context.Context, metric *schemas.Metric) error
	SaveAll(ctx context.Context, metrics []schemas.Metric) error
	Ping(ctx context.Context) error
	Close() error
}
