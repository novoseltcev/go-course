package storages

import (
	"context"
	"sort"

	"github.com/novoseltcev/go-course/internal/schemas"
)

type MemStorage struct {
	Data map[string]map[string]schemas.MetricValues
}

func NewMemStorage() *MemStorage {
	m := make(map[string]map[string]schemas.MetricValues)
	m[schemas.Counter] = make(map[string]schemas.MetricValues)
	m[schemas.Gauge] = make(map[string]schemas.MetricValues)

	return &MemStorage{m}
}

func (s MemStorage) GetOne(_ context.Context, id, mType string) (*schemas.Metric, error) {
	result, ok := s.Data[mType][id]
	if !ok {
		return nil, ErrNotFound
	}

	return &schemas.Metric{ID: id, MType: mType, Value: result.Value, Delta: result.Delta}, nil
}

func (s *MemStorage) GetAll(_ context.Context) ([]schemas.Metric, error) {
	result := make([]schemas.Metric, 0)

	for Type := range s.Data {
		data := s.Data[Type]

		names := make([]string, 0, len(data))
		for name := range data {
			names = append(names, name)
		}

		sort.Strings(names)

		for _, name := range names {
			d := data[name]

			result = append(result, schemas.Metric{ID: name, MType: Type, Value: d.Value, Delta: d.Delta})
		}
	}

	return result, nil
}

func (s *MemStorage) Save(_ context.Context, metric *schemas.Metric) error {
	stored, ok := s.Data[metric.MType][metric.ID]
	if ok && metric.MType == "counter" && stored.Delta != nil {
		*stored.Delta += *metric.Delta
	} else {
		s.Data[metric.MType][metric.ID] = schemas.MetricValues{Value: metric.Value, Delta: metric.Delta}
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
