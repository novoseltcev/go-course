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


func GetOneMetricFromJSON(storage storage.MetricStorager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var metric schema.Metrics
        if err := json.UnmarshalFromReader(r.Body, &metric); err != nil {
			log.Warn(err.Error())
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

		w.Header().Set("Content-Type", "application/json")
		switch metric.MType {
		case "gauge":
			result, err := utils.RetryPgSelect(ctx, func() (*model.Metric, error) {
				return storage.GetByName(ctx, metric.ID, metric.MType)
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if result == nil {
				http.Error(w, "Metric not found", http.StatusNotFound)
				return
			}

			metric.Value = result.Value
		case "counter":
			result, err := utils.RetryPgSelect(ctx, func() (*model.Metric, error) {
				return storage.GetByName(ctx, metric.ID, metric.MType)
			})

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if result == nil {
				http.Error(w, "Metric not found", http.StatusNotFound)
				return
			}
			
			metric.Delta = result.Delta
		default:
			log.Warn(http.StatusText(http.StatusBadRequest))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if _, err := json.MarshalToWriter(metric, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		w.WriteHeader(http.StatusOK)
	}
}
