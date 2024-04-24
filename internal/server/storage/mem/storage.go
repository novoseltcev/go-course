package mem

import (
	"context"
	"sort"

	"github.com/novoseltcev/go-course/internal/model"
)

type Storage struct {
	Metrics map[string]map[string]model.Metric
}

func (s Storage) GetByName(ctx context.Context, name, Type string) (*model.Metric, error) {
	result, ok := s.Metrics[Type][name]
	if !ok {
		return nil, nil
	}

	return &result, nil
}

func (s *Storage) GetAll(ctx context.Context) ([]model.Metric, error) {
	result := make([]model.Metric, 0)

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

func (s *Storage) Save(ctx context.Context, metric model.Metric) error {
	saved, ok := s.Metrics[metric.Type][metric.Name]
	if metric.Type == "counter" && ok && saved.Delta != nil {
		*saved.Delta += *metric.Delta
	} else {
		s.Metrics[metric.Type][metric.Name] = metric
	}
	return nil
}

func (s *Storage) SaveAll(ctx context.Context, metrics []model.Metric) error {
	for _, metric := range metrics {
		s.Save(ctx, metric)
	}
	return nil
}
