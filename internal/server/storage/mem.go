package storage

import (
	"sort"

	"github.com/novoseltcev/go-course/internal/types"
)

type MemStorage[T types.Counter | types.Gauge] struct {
	Metrics map[string]T
}

func (s *MemStorage[T]) GetAll() []types.Metric[T] {
	names := make([]string, 0, len(s.Metrics))
	for name := range s.Metrics {
		names = append(names, name)
	}
	sort.Strings(names)

	result := make([]types.Metric[T], 0, len(s.Metrics))
	for _, name := range names {
		result = append(result, types.Metric[T]{Name: name, Value: s.Metrics[name]})
	}

	return result
}

func (s *MemStorage[T]) Update(name string, value T) {
	switch any(value).(type) {
	case types.Counter:
		s.Metrics[name] = s.Metrics[name] + value
	case types.Gauge:
		s.Metrics[name] = value
	}
}

func (s MemStorage[T]) GetByName(name string) *T {
	val, ok := s.Metrics[name]
	if !ok {
		return nil
	}

	return &val
}
