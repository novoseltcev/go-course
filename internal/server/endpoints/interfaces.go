package endpoints

import "github.com/novoseltcev/go-course/internal/types"


type MetricStorager[T types.Counter | types.Gauge] interface {
	GetByName(string) *T
	GetAll() []types.Metric[T]
	Update(string, T)
}
