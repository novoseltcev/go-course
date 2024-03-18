package endpoints

import (
	"net/http"
	"strconv"
	"github.com/go-chi/chi/v5"
	"github.com/novoseltcev/go-course/internal/server/storage"
)


func UpdateMetric(counterStorage *storage.Storage[storage.Counter], gaugeStorage *storage.Storage[storage.Gauge]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")
		metricValue := chi.URLParam(r, "metricValue")

		switch metricType {
		case "gauge":
			value, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			
			(*gaugeStorage).Update(metricName, storage.Gauge(value))
		case "counter":
			value, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			(*counterStorage).Update(metricName, storage.Counter(value))
		default:
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
