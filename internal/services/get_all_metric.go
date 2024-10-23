package services

import (
	"context"

	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/internal/storages"
	"github.com/novoseltcev/go-course/pkg/retry"
)

func GetAllMetric(
	ctx context.Context,
	storage storages.MetricStorager,
) ([]schemas.Metric, error) {
	result, err := retry.PgSelect(ctx, func() ([]schemas.Metric, error) {
		return storage.GetAll(ctx)
	}, nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}
