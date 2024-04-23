package endpoints

import (
	"net/http"

	json "github.com/mailru/easyjson"
	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/model"
	"github.com/novoseltcev/go-course/internal/schema"
	"github.com/novoseltcev/go-course/internal/server/storage"
)


func UpdateMetricFromJSON(counterStorage *storage.MetricStorager[model.Counter], gaugeStorage *storage.MetricStorager[model.Gauge]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

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
			(*gaugeStorage).Update(ctx, metric.ID, model.Gauge(*metric.Value))
		case "counter":
			if metric.Delta == nil {
				log.Error("counter metric has nil delta")
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
            	return
			}
			(*counterStorage).Update(ctx, metric.ID, model.Counter(*metric.Delta))
		default:
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
