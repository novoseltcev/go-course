package storage

type Metric[T int64 | float64] struct {
	Name string
	Value T
}

type Storage[T int64 | float64] interface {
	GetAll() []Metric[T]
	Update(string, T)
}
