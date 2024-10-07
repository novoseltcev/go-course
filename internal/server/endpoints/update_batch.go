package endpoints

import (
	"net/http"

	json "github.com/mailru/easyjson"
	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/schema"
	"github.com/novoseltcev/go-course/internal/server/services"
	"github.com/novoseltcev/go-course/internal/server/storage"
)

func UpdateMetricsBatch(storage storage.MetricStorager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var metrics schema.MetricSlice
		if err := json.UnmarshalFromReader(r.Body, &metrics); err != nil {
			log.WithError(err).Error("unmarshalable body")
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		batch := make([]schema.Metric, 0, len(metrics))

		for _, metric := range metrics {
			switch metric.MType {
			case schema.Gauge:
				if metric.Value == nil {
					log.Error("gauge metric has nil value")

					continue
				}
			case schema.Counter:
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
			err := services.SaveMetricsBatch(ctx, storage, pgRetries, batch)
			if err != nil {
				log.WithError(err).Error("cannot save metrics")
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}
