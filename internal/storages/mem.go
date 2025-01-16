package storages

import (
	"context"
	"sort"

	"github.com/novoseltcev/go-course/internal/schemas"
)

type MemStorage struct {
	m map[string]map[string]schemas.Metric
}

func NewMemStorage() *MemStorage {
	m := make(map[string]map[string]schemas.Metric)
	m[schemas.Counter] = make(map[string]schemas.Metric)
	m[schemas.Gauge] = make(map[string]schemas.Metric)

	return &MemStorage{m}
}

func (s MemStorage) GetOne(_ context.Context, id, mType string) (*schemas.Metric, error) {
	result, ok := s.m[mType][id]
	if !ok {
		return nil, ErrNotFound
	}

	return &result, nil
}

func (s *MemStorage) GetAll(_ context.Context) ([]schemas.Metric, error) {
	result := make([]schemas.Metric, 0)

	for Type := range s.m {
		data := s.m[Type]

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
	saved, ok := s.m[metric.MType][metric.ID]
	if metric.MType == "counter" && ok && saved.Delta != nil {
		*saved.Delta += *metric.Delta
	} else {
		s.m[metric.MType][metric.ID] = *metric
	}

	return nil
}

func (s *MemStorage) SaveBatch(ctx context.Context, metrics []schemas.Metric) error {
	for _, metric := range metrics {
		if err := s.Save(ctx, &metric); err != nil {
			return err
		}
	}

	return nil
}
