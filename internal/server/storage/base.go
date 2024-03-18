package storage

type Counter int64
type Gauge float64

type Metric[T Counter | Gauge] struct {
	Name string
	Value T
}

type Storage[T Counter | Gauge] interface {
	GetByName(string) *T
	GetAll() []Metric[T]
	Update(string, T)
}
