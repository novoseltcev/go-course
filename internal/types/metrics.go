package types

type Counter int64
type Gauge float64

type Metric[T Counter | Gauge] struct {
	Name string
	Value T
}
