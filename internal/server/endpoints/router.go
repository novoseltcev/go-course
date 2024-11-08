package endpoints

import (
	"github.com/go-chi/chi/v5"

	"github.com/novoseltcev/go-course/internal/storages"
)

func NewAPIRouter(storage storages.MetricStorager) *chi.Mux {
	r := chi.NewRouter()

	r.Get(`/ping`, Ping(storage))
	r.Get(`/`, Index(storage))
	r.Post(`/update/{type}/{id}/{value}`, UpdateMetric(storage))
	r.Get(`/value/{type}/{id}`, GetOneMetric(storage))
	r.Post(`/update/`, UpdateMetricFromJSON(storage))
	r.Post(`/value/`, GetOneMetricFromJSON(storage))
	r.Post(`/updates/`, UpdateMetricsBatch(storage))

	return r
}
