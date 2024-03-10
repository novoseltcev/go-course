package storage


type MemStorage[T int64 | float64] struct {
	metrics map[string]T
}

func (storage MemStorage[T]) GetAll() []Metric[T] {
	var result = make([]Metric[T], len(storage.metrics))
	for name, value := range storage.metrics {
		result = append(result, Metric[T]{Name: name, Value: value})
	}
	return result
}

func (storage MemStorage[T]) Update(name string, value T) {
	switch any(value).(type) {
	case float64:
		storage.metrics[name] = value
	case int64:
		storage.metrics[name] = storage.metrics[name] + value
	}
}
