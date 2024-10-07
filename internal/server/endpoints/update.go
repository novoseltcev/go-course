package endpoints

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/novoseltcev/go-course/internal/schema"
	"github.com/novoseltcev/go-course/internal/server/services"
	"github.com/novoseltcev/go-course/internal/server/storage"
)

func UpdateMetric(storage storage.MetricStorager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")
		metricValue := chi.URLParam(r, "metricValue")

		metric := schema.Metric{ID: metricName, MType: metricType}

		if metricType == schema.Gauge {
			value, err := strconv.ParseFloat(metricValue, 64)
			if err == nil {
				metric.Value = &value
			}
		}

		if metricType == schema.Counter {
			value, err := strconv.ParseInt(metricValue, 10, 64)
			if err == nil {
				metric.Delta = &value
			}
		}

		if err := services.SaveMetric(ctx, storage, pgRetries, &metric); err != nil {
			var statusCode int

			switch {
			case errors.Is(err, services.ErrInvalidType):
				statusCode = http.StatusBadRequest
			case errors.Is(err, services.ErrInvalidValue), errors.Is(err, services.ErrInvalidDelta):
				statusCode = http.StatusBadRequest
			default:
				statusCode = http.StatusInternalServerError
			}

			http.Error(w, err.Error(), statusCode)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
