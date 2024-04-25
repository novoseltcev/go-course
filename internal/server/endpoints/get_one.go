package endpoints

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/novoseltcev/go-course/internal/model"
	"github.com/novoseltcev/go-course/internal/server/storage"
	"github.com/novoseltcev/go-course/internal/utils"
)


func GetOneMetric(storage *storage.MetricStorager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")

		metricValue := ""
		switch metricType {
		case "gauge":
			result, err := utils.RetryPgSelect(ctx, func() (*model.Metric, error) {
				return (*storage).GetByName(ctx, metricName, metricType)
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if result != nil {
				metricValue = strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", *result.Value), "0"), ".")
			}
		case "counter":
			result, err := utils.RetryPgSelect(ctx, func() (*model.Metric, error) {
				return (*storage).GetByName(ctx, metricName, metricType)
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if result != nil {
				metricValue = strconv.Itoa(int(*result.Delta))
			}

		default:
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		
		if metricValue == "" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		io.WriteString(w, metricValue)
	}
}
