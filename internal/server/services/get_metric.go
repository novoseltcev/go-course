package services

import (
	"context"
	"errors"

	"github.com/novoseltcev/go-course/internal/schema"
	"github.com/novoseltcev/go-course/internal/server/storage"
	"github.com/novoseltcev/go-course/internal/utils"
)

var (
	ErrMetricNotFound = errors.New("metric not found")
	ErrInvalidType    = errors.New("invalid metric type")
)

func GetMetric(
	ctx context.Context,
	storage storage.MetricStorager,
	pgRetries uint,
	metricName, metricType string,
) (*schema.Metric, error) {
	if metricType != schema.Gauge && metricType != schema.Counter {
		return nil, ErrInvalidType
	}

	result, err := utils.RetryPgSelect(ctx, func() (*schema.Metric, error) {
		return storage.GetByName(ctx, metricName, metricType)
	}, pgRetries)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, ErrMetricNotFound
	}

	return result, nil
}
