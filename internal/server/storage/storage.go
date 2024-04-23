package storage

import (
	"context"

	"github.com/novoseltcev/go-course/internal/model"
)


type MetricStorager[T model.Counter | model.Gauge] interface {
	GetByName(context.Context, string) *T
	GetAll(context.Context) []model.Metric[T]
	Update(context.Context, string, T)
}
