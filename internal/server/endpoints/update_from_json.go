package endpoints

import (
	"net/http"

	json "github.com/mailru/easyjson"
	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/model"
	"github.com/novoseltcev/go-course/internal/schema"
)


func UpdateMetricFromJSON(counterStorage *MetricStorager[model.Counter], gaugeStorage *MetricStorager[model.Gauge]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metric schema.Metrics
        if err := json.UnmarshalFromReader(r.Body, &metric); err != nil {
			log.Error(err.Error())
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

		switch metric.MType {
		case "gauge":
			if metric.Value == nil {
				log.Error("gauge metric has nil value")
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			(*gaugeStorage).Update(metric.ID, model.Gauge(*metric.Value))
		case "counter":
			if metric.Delta == nil {
				log.Error("counter metric has nil delta")
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
            	return
			}
			(*counterStorage).Update(metric.ID, model.Counter(*metric.Delta))
		default:
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
