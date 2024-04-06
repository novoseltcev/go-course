package endpoints

import (
	"net/http"

	json "github.com/mailru/easyjson"

	"github.com/novoseltcev/go-course/internal/model"
	"github.com/novoseltcev/go-course/internal/schema"
)


func GetOneMetricFromJSON(counterStorage *MetricStorager[model.Counter], gaugeStorage *MetricStorager[model.Gauge]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metric schema.Metrics
        if err := json.UnmarshalFromReader(r.Body, &metric); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

		w.Header().Set("Content-Type", "application/json")
		switch metric.MType {
		case "gauge":
			result := (*gaugeStorage).GetByName(metric.ID)
			var value float64
			if result == nil {
				value = 0
			} else {
				value = float64(*result)
			}
			metric.Value = &value
		case "counter":
			result := (*counterStorage).GetByName(metric.ID)

			var value int64
			if result == nil {
				value = 0
			} else {
				value = int64(*result)
			}
			metric.Delta = &value
		default:
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
