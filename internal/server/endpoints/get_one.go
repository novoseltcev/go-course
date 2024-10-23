package endpoints

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/internal/services"
	"github.com/novoseltcev/go-course/internal/storages"
)

func GetOneMetric(storage storages.MetricStorager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")

		result, err := services.GetMetric(ctx, storage, metricName, metricType)
		if err != nil {
			var statusCode int

			switch {
			case errors.Is(err, services.ErrInvalidType):
				statusCode = http.StatusBadRequest
			case errors.Is(err, services.ErrMetricNotFound):
				statusCode = http.StatusNotFound
			default:
				statusCode = http.StatusInternalServerError
			}

			http.Error(w, err.Error(), statusCode)

			return
		}

		var metricValue string

		switch metricType {
		case schemas.Gauge:
			metricValue = strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", *result.Value), "0"), ".")
		case schemas.Counter:
			metricValue = strconv.Itoa(int(*result.Delta))
		default:
			panic("unreachable")
		}

		if _, err := io.WriteString(w, metricValue); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
