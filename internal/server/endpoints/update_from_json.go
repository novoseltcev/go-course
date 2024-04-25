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


func UpdateMetricFromJSON(storage *storage.MetricStorager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var s schema.Metrics
        if err := json.UnmarshalFromReader(r.Body, &s); err != nil {
			log.Error(err.Error())
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

		var metric model.Metric
		switch s.MType {
		case "gauge":
			if s.Value == nil {
				log.Error("gauge metric has nil value")
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			metric = model.Metric{Name: s.ID, Type: s.MType, Value: s.Value}
		case "counter":
			if s.Delta == nil {
				log.Error("counter metric has nil delta")
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
            	return
			}
			
			metric = model.Metric{Name: s.ID, Type: s.MType, Delta: s.Delta}
		default:
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		
		err := utils.RetryPgExec(ctx, func() error {
			return (*storage).Save(ctx, metric)
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
