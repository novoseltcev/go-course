package storages

import (
	"context"
	"sort"

	"github.com/novoseltcev/go-course/internal/schemas"
)

type MemStorage struct {
	Metrics map[string]map[string]schemas.Metric
}

func NewMemStorage() *MemStorage {
	metrics := make(map[string]map[string]schemas.Metric)
	metrics["counter"] = make(map[string]schemas.Metric)
	metrics["gauge"] = make(map[string]schemas.Metric)

	return &MemStorage{Metrics: metrics}
}

func (s MemStorage) GetByName(_ context.Context, name, metricType string) (*schemas.Metric, error) {
	result, ok := s.Metrics[metricType][name]
	if !ok {
		return nil, nil //nolint:nilnil
	}

	return &result, nil
}

func (s *MemStorage) GetAll(_ context.Context) ([]schemas.Metric, error) {
	result := make([]schemas.Metric, 0)

	for Type := range s.Metrics {
		data := s.Metrics[Type]

		names := make([]string, 0, len(data))
		for name := range data {
			names = append(names, name)
		}

		sort.Strings(names)

		for _, name := range names {
			result = append(result, data[name])
		}
	}

	return result, nil
}

func (s *MemStorage) Save(_ context.Context, metric *schemas.Metric) error {
	saved, ok := s.Metrics[metric.MType][metric.ID]
	if metric.MType == "counter" && ok && saved.Delta != nil {
		*saved.Delta += *metric.Delta
	} else {
		s.Metrics[metric.MType][metric.ID] = *metric
	}

	return nil
}

func (s *MemStorage) SaveAll(ctx context.Context, metrics []schemas.Metric) error {
	for _, metric := range metrics {
		if err := s.Save(ctx, &metric); err != nil {
			return err
		}
	}

	return nil
}

func (s *MemStorage) Ping(_ context.Context) error {
	return nil
}

func (s *MemStorage) Close() error {
	return nil
}
