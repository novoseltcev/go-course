package services

import (
	"context"

	"github.com/novoseltcev/go-course/internal/schema"
	"github.com/novoseltcev/go-course/internal/server/storage"
	"github.com/novoseltcev/go-course/internal/utils"
)

func GetAllMetric(
	ctx context.Context,
	storage storage.MetricStorager,
	pgRetries uint,
) ([]schema.Metric, error) {
	result, err := utils.RetryPgSelect(ctx, func() ([]schema.Metric, error) {
		return storage.GetAll(ctx)
	}, pgRetries)
	if err != nil {
		return nil, err
	}

	return result, nil
}
