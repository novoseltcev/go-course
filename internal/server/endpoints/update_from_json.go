package endpoints

import (
	"net/http"

	json "github.com/mailru/easyjson"
	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/internal/storages"
)

func UpdateMetricFromJSON(storager storages.MetricStorager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var metric schemas.Metric
		if err := json.UnmarshalFromReader(r.Body, &metric); err != nil {
			log.WithError(err).Error("unmarshalable body")
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		if err := metric.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		if err := storager.Save(ctx, &metric); err != nil {
			log.WithError(err).Error("failed to save metric")
			http.Error(w, "failed to save metric", http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
