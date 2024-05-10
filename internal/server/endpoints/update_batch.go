package endpoints

import (
	"net/http"

	json "github.com/mailru/easyjson"
	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/model"
	"github.com/novoseltcev/go-course/internal/schema"
	"github.com/novoseltcev/go-course/internal/server/storage"
	"github.com/novoseltcev/go-course/internal/utils"
)


func UpdateMetricsBatch(storage storage.MetricStorager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var metrics schema.MetricsSlice
        if err := json.UnmarshalFromReader(r.Body, &metrics); err != nil {
			log.WithError(err).Error("unmarshalable body")
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

		batch := make([]model.Metric, 0, len(metrics))

		for _, metric := range metrics {
			switch metric.MType {
			case "gauge":
				if metric.Value == nil {
					log.Error("gauge metric has nil value")
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}

				batch = append(batch, model.Metric{Name: metric.ID, Type: metric.MType, Value: metric.Value})
			case "counter":
				if metric.Delta == nil {
					log.Error("counter metric has nil delta")
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}

				batch = append(batch, model.Metric{Name: metric.ID, Type: metric.MType, Delta: metric.Delta})
			default:
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
		}

		if len(batch) != 0 {
			err := utils.RetryPgExec(ctx, func() error {
				return storage.SaveAll(ctx, batch)
			})
			if err != nil {
				log.WithError(err).Error("cannot save metrics")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		
		w.WriteHeader(http.StatusOK)
	}
}
