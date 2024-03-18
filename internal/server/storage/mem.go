package storage

import "sort"

type MemStorage[T Counter | Gauge] struct {
	Metrics map[string]T
}

func (s MemStorage[T]) GetAll() []Metric[T] {
	names := make([]string, 0, len(s.Metrics))
	for name := range s.Metrics {
		names = append(names, name)
	}
	sort.Strings(names)

	result := make([]Metric[T], 0, len(s.Metrics))
	for _, name := range names {
		result = append(result, Metric[T]{Name: name, Value: s.Metrics[name]})
	}

	return result
}

func (s MemStorage[T]) Update(name string, value T) {
	switch any(value).(type) {
	case Counter:
		s.Metrics[name] = s.Metrics[name] + value
	case Gauge:
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
