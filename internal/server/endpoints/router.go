package endpoints

import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/novoseltcev/go-course/internal/server/storage"
)


func GetRouter(counterStorage *storage.Storage[storage.Counter], gaugeStorage *storage.Storage[storage.Gauge]) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	
	r.Post(`/update/{metricType}/{metricName}/{metricValue}`, UpdateMetric(counterStorage, gaugeStorage))
	r.Get(`/value/{metricType}/{metricName}`, GetOneMetric(counterStorage, gaugeStorage))
	r.Get(`/`, Index(counterStorage, gaugeStorage))
	return r
}
