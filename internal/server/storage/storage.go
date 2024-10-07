package storage

import (
	"context"

	"github.com/novoseltcev/go-course/internal/schema"
)

type MetricStorager interface {
	GetByName(ctx context.Context, name, Type string) (*schema.Metric, error)
	GetAll(ctx context.Context) ([]schema.Metric, error)
	Save(ctx context.Context, metric *schema.Metric) error
	SaveAll(ctx context.Context, metrics []schema.Metric) error
	Ping(ctx context.Context) error
	Close() error
}
