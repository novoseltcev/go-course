package endpoints

import "github.com/novoseltcev/go-course/internal/model"


type MetricStorager[T model.Counter | model.Gauge] interface {
	GetByName(string) T
	GetAll() []model.Metric[T]
	Update(string, T)
}
