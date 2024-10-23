package services

import (
	"context"
	"errors"

	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/internal/storages"
	"github.com/novoseltcev/go-course/pkg/retry"
)

var (
	ErrMetricNotFound = errors.New("metric not found")
	ErrInvalidType    = errors.New("invalid metric type")
)

func GetMetric(
	ctx context.Context,
	storage storages.MetricStorager,
	metricName, metricType string,
) (*schemas.Metric, error) {
	if metricType != schemas.Gauge && metricType != schemas.Counter {
		return nil, ErrInvalidType
	}

	result, err := retry.PgSelect(ctx, func() (*schemas.Metric, error) {
		return storage.GetByName(ctx, metricName, metricType)
	}, nil)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, ErrMetricNotFound
	}

	return result, nil
}
