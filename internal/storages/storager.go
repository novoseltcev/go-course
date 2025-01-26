package storages

import (
	"context"
	"errors"

	"github.com/novoseltcev/go-course/internal/schemas"
)

//go:generate mockgen -source=storager.go  -destination=../../mocks/storager_mock.go -package=mocks -typed

var ErrNotFound = errors.New("metric not found")

// Storager is an interface for metric storage.
//
// It provides methods for getting and saving metrics.
// The interface also includes a Ping method for testing the connection and a Close method for closing the connection.
//
// TODO: Separate interface on many interfaces.
type MetricStorager interface {
	GetOne(ctx context.Context, id, mType string) (*schemas.Metric, error)
	GetAll(ctx context.Context) ([]schemas.Metric, error)
	Save(ctx context.Context, metric *schemas.Metric) error
	SaveBatch(ctx context.Context, metrics []schemas.Metric) error
}
