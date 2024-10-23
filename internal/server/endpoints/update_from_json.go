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

func UpdateMetricFromJSON(storage storages.MetricStorager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var reqBody schemas.Metric
		if err := json.UnmarshalFromReader(r.Body, &reqBody); err != nil {
			log.WithError(err).Error("unmarshalable body")
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		if err := services.SaveMetric(ctx, storage, &reqBody); err != nil {
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
