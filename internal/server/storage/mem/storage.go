package mem

import (
	"context"
	"sort"

	"github.com/novoseltcev/go-course/internal/schema"
	s "github.com/novoseltcev/go-course/internal/server/storage"
)

type storage struct {
	Metrics map[string]map[string]schema.Metric
}

func New() s.MetricStorager {
	metrics := make(map[string]map[string]schema.Metric)
	metrics["counter"] = make(map[string]schema.Metric)
	metrics["gauge"] = make(map[string]schema.Metric)

	return &storage{Metrics: metrics}
}

func (s storage) GetByName(_ context.Context, name, metricType string) (*schema.Metric, error) {
	result, ok := s.Metrics[metricType][name]
	if !ok {
		return nil, nil //nolint:nilnil
	}

	return &result, nil
}

func (s *storage) GetAll(_ context.Context) ([]schema.Metric, error) {
	result := make([]schema.Metric, 0)

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

func (s *storage) Save(_ context.Context, metric *schema.Metric) error {
	saved, ok := s.Metrics[metric.MType][metric.ID]
	if metric.MType == "counter" && ok && saved.Delta != nil {
		*saved.Delta += *metric.Delta
	} else {
		s.Metrics[metric.MType][metric.ID] = *metric
	}

	return nil
}

func (s *storage) SaveAll(ctx context.Context, metrics []schema.Metric) error {
	for _, metric := range metrics {
		if err := s.Save(ctx, &metric); err != nil {
			return err
		}
	}

	return nil
}

func (s *storage) Ping(_ context.Context) error {
	return nil
}

func (s *storage) Close() error {
	return nil
}
