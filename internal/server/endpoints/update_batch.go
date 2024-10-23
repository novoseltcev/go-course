package endpoints

import (
	"net/http"

	json "github.com/mailru/easyjson"
	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/internal/services"
	"github.com/novoseltcev/go-course/internal/storages"
)

func UpdateMetricsBatch(storage storages.MetricStorager) http.HandlerFunc {
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
			switch metric.MType {
			case schemas.Gauge:
				if metric.Value == nil {
					log.Error("gauge metric has nil value")

					continue
				}
			case schemas.Counter:
				if metric.Delta == nil {
					log.Error("counter metric has nil delta")

					continue
				}
			default:
				log.WithField("type", metric.MType).Error("invalid metric type")

				continue
			}

			batch = append(batch, metric)
		}

		if len(batch) != 0 {
			err := services.SaveMetricsBatch(ctx, storage, batch)
			if err != nil {
				log.WithError(err).Error("cannot save metrics")
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}
