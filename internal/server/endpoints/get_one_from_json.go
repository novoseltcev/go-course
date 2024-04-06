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
			result := float64((*gaugeStorage).GetByName(metric.ID))
			metric.Value = &result
		case "counter":
			result := int64((*counterStorage).GetByName(metric.ID))
			metric.Delta = &result
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
