package endpoints

import (
	"net/http"

	json "github.com/mailru/easyjson"
	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/model"
	"github.com/novoseltcev/go-course/internal/schema"
	"github.com/novoseltcev/go-course/internal/server/storage"
)


func GetOneMetricFromJSON(counterStorage *storage.MetricStorager[model.Counter], gaugeStorage *storage.MetricStorager[model.Gauge]) http.HandlerFunc {
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
			result := (*gaugeStorage).GetByName(ctx, metric.ID)
			if result == nil {
				http.Error(w, "Metric not found", http.StatusNotFound)
				return
			}

			metric.Value = (*float64)(result)
		case "counter":
			result := (*counterStorage).GetByName(ctx, metric.ID)
			if result == nil {
				http.Error(w, "Metric not found", http.StatusNotFound)
				return
			}
			
			metric.Delta = (*int64)(result)
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
