package mem

import (
	"sort"

	"github.com/novoseltcev/go-course/internal/model"
)

type Storage[T model.Counter | model.Gauge] struct {
	Metrics map[string]T
}

func (s *Storage[T]) GetAll() []model.Metric[T] {
	names := make([]string, 0, len(s.Metrics))
	for name := range s.Metrics {
		names = append(names, name)
	}
	sort.Strings(names)

	result := make([]model.Metric[T], 0, len(s.Metrics))
	for _, name := range names {
		result = append(result, model.Metric[T]{Name: name, Value: s.Metrics[name]})
	}

	return result
}

func (s *Storage[T]) Update(name string, value T) {
	switch any(value).(type) {
	case model.Counter:
		s.Metrics[name] = s.Metrics[name] + value
	case model.Gauge:
		s.Metrics[name] = value
	}
}

func (s Storage[T]) GetByName(name string) *T {
	result, ok := s.Metrics[name]
	if !ok {
		return nil
	}

	return &result
}
