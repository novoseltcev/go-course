package endpoints

import (
	"net/http"

	json "github.com/mailru/easyjson"
	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/model"
	"github.com/novoseltcev/go-course/internal/schema"
	"github.com/novoseltcev/go-course/internal/server/storage"
)


func UpdateMetricsBatch(storage *storage.MetricStorager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var metrics schema.MetricsList
        if err := json.UnmarshalFromReader(r.Body, &metrics); err != nil {
			log.Error(err.Error())
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

		batch := make([]model.Metric, 0, len(metrics.Metrics))

		for _, metric := range metrics.Metrics {
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
			if err := (*storage).SaveAll(ctx, batch); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		
		w.WriteHeader(http.StatusOK)
	}
}
