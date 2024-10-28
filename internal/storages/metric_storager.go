package storages

import (
	"context"

	"github.com/novoseltcev/go-course/internal/schemas"
)

// Storager is an interface for metric storage.
//
// It provides methods for getting and saving metrics.
// The interface also includes a Ping method for testing the connection and a Close method for closing the connection.
//
// TODO: Separate interface on many interfaces.
type MetricStorager interface {
	GetByName(ctx context.Context, name, Type string) (*schemas.Metric, error)
	GetAll(ctx context.Context) ([]schemas.Metric, error)
	Save(ctx context.Context, metric *schemas.Metric) error
	SaveAll(ctx context.Context, metrics []schemas.Metric) error
	Ping(ctx context.Context) error
	Close() error
}
