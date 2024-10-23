package services

import (
	"context"
	"errors"

	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/internal/storages"
	"github.com/novoseltcev/go-course/pkg/retry"
)

var (
	ErrInvalidValue = errors.New("invalid gauge value")
	ErrInvalidDelta = errors.New("invalid counter value")
)

func SaveMetric(
	ctx context.Context,
	storage storages.MetricStorager,
	metric *schemas.Metric,
) error {
	switch metric.MType {
	case schemas.Gauge:
		if metric.Value == nil {
			return ErrInvalidValue
		}

	case schemas.Counter:
		if metric.Delta == nil {
			return ErrInvalidDelta
		}

	default:
		return ErrInvalidType
	}

	return retry.PgExec(ctx, func() error {
		return storage.Save(ctx, metric)
	}, nil)
}

func SaveMetricsBatch(
	ctx context.Context,
	storage storages.MetricStorager,
	metrics []schemas.Metric,
) error {
	return retry.PgExec(ctx, func() error {
		return storage.SaveAll(ctx, metrics)
	}, nil)
}
