package storage

type MemStorage[T Counter | Gauge] struct {
	Metrics map[string]T
}

func (storage MemStorage[T]) GetAll() []Metric[T] {
	var result = make([]Metric[T], len(storage.Metrics))
	for name, value := range storage.Metrics {
		result = append(result, Metric[T]{Name: name, Value: value})
	}
	return result
}

func (storage MemStorage[T]) Update(name string, value T) {
	switch any(value).(type) {
	case Counter:
		storage.Metrics[name] = storage.Metrics[name] + value
	case Gauge:
		storage.Metrics[name] = value
	}
}
