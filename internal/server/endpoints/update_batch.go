package endpoints

import (
	"net/http"

	json "github.com/mailru/easyjson"
	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/internal/storages"
)

func UpdateMetricsBatch(storager storages.MetricStorager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var metrics schemas.MetricSlice
		if err := json.UnmarshalFromReader(r.Body, &metrics); err != nil {
			log.WithError(err).Error("unmarshalable body")
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		batch := make([]schemas.Metric, 0, len(metrics))

		for _, metric := range metrics {
			if err := metric.Validate(); err != nil {
				log.WithError(err).Error("invalid metric")

				continue
			}

			batch = append(batch, metric)
		}

		if len(batch) != 0 {
			err := storager.SaveBatch(ctx, batch)
			if err != nil {
				log.WithError(err).Error("failed to save metrics")
				http.Error(w, "failed to save metrics", http.StatusInternalServerError)

				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}
