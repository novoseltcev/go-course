package endpoints

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/novoseltcev/go-course/internal/model"
	"github.com/novoseltcev/go-course/internal/server/middlewares"
)


func GetRouter(counterStorage *MetricStorager[model.Counter], gaugeStorage *MetricStorager[model.Gauge]) http.Handler {
	r := chi.NewRouter()
	r.Use(middlewares.Gzip, middlewares.Logger)

	r.Post(`/update/{metricType}/{metricName}/{metricValue}`, UpdateMetric(counterStorage, gaugeStorage))
	r.Get(`/value/{metricType}/{metricName}`, GetOneMetric(counterStorage, gaugeStorage))
	r.Post(`/update/`, UpdateMetricFromJSON(counterStorage, gaugeStorage))
	r.Post(`/value/`, GetOneMetricFromJSON(counterStorage, gaugeStorage))
	r.Get(`/`, Index(counterStorage, gaugeStorage))
	return r
}
