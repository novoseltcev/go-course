package endpoints

import (
	"errors"
	"net/http"

	json "github.com/mailru/easyjson"
	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/internal/storages"
)

func GetOneMetricFromJSON(storager storages.MetricStorager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var metric schemas.MetricIdentifier
		if err := json.UnmarshalFromReader(r.Body, &metric); err != nil {
			log.Warn(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		if err := metric.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		result, err := storager.GetOne(ctx, metric.ID, metric.MType)
		if err != nil {
			if errors.Is(err, storages.ErrNotFound) {
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				log.WithField("id", metric.ID).WithField("type", metric.MType).WithError(err).Error("failed to get metric")
				http.Error(w, "failed to get metric", http.StatusInternalServerError)
			}

			return
		}

		if _, err := json.MarshalToWriter(result, w); err != nil {
			log.WithError(err).Error("failed to serialize metric")
			http.Error(w, "failed to serialize metric", http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
