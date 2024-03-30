package endpoints

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/novoseltcev/go-course/internal/types"
	"github.com/novoseltcev/go-course/internal/server/middlewares"
)


func GetRouter(counterStorage *MetricStorager[types.Counter], gaugeStorage *MetricStorager[types.Gauge]) http.Handler {
	r := chi.NewRouter()
	r.Use(middlewares.Logger)
	
	r.Post(`/update/{metricType}/{metricName}/{metricValue}`, UpdateMetric(counterStorage, gaugeStorage))
	r.Get(`/value/{metricType}/{metricName}`, GetOneMetric(counterStorage, gaugeStorage))
	r.Get(`/`, Index(counterStorage, gaugeStorage))
	return r
}
