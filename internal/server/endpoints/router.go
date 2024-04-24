package endpoints

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/novoseltcev/go-course/internal/server/middlewares"
	"github.com/novoseltcev/go-course/internal/server/storage"
)


func GetRouter(db *sqlx.DB, storage *storage.MetricStorager) http.Handler {
	r := chi.NewRouter()
	r.Use(middlewares.Gzip, middlewares.Logger)

	r.Post(`/updates/`, UpdateMetricsBatch(storage))
	r.Post(`/update/{metricType}/{metricName}/{metricValue}`, UpdateMetric(storage))
	r.Get(`/value/{metricType}/{metricName}`, GetOneMetric(storage))
	r.Post(`/update/`, UpdateMetricFromJSON(storage))
	r.Post(`/value/`, GetOneMetricFromJSON(storage))
	r.Get(`/ping`, Ping(db))
	r.Get(`/`, Index(storage))
	return r
}
