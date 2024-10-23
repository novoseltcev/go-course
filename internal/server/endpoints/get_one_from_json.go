package endpoints

import (
	"errors"
	"net/http"

	json "github.com/mailru/easyjson"
	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/internal/services"
	"github.com/novoseltcev/go-course/internal/storages"
)

func GetOneMetricFromJSON(storage storages.MetricStorager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var metric schemas.Metric
		if err := json.UnmarshalFromReader(r.Body, &metric); err != nil {
			log.Warn(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		w.Header().Set("Content-Type", "application/json")

		result, err := services.GetMetric(ctx, storage, metric.ID, metric.MType)
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

		if _, err := json.MarshalToWriter(result, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
