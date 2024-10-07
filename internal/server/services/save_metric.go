package services

import (
	"context"
	"errors"

	"github.com/novoseltcev/go-course/internal/schema"
	"github.com/novoseltcev/go-course/internal/server/storage"
	"github.com/novoseltcev/go-course/internal/utils"
)

var (
	ErrInvalidValue = errors.New("invalid gauge value")
	ErrInvalidDelta = errors.New("invalid counter value")
)

func SaveMetric(
	ctx context.Context,
	storage storage.MetricStorager,
	pgRetries uint,
	metric *schema.Metric,
) error {
	switch metric.MType {
	case schema.Gauge:
		if metric.Value == nil {
			return ErrInvalidValue
		}

	case schema.Counter:
		if metric.Delta == nil {
			return ErrInvalidDelta
		}

	default:
		return ErrInvalidType
	}

	return utils.RetryPgExec(ctx, func() error {
		return storage.Save(ctx, metric)
	}, pgRetries)
}

func SaveMetricsBatch(
	ctx context.Context,
	storage storage.MetricStorager,
	pgRetries uint,
	metrics []schema.Metric,
) error {
	return utils.RetryPgExec(ctx, func() error {
		return storage.SaveAll(ctx, metrics)
	}, pgRetries)
}
