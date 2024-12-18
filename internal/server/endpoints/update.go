package endpoints

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/internal/storages"
)

func UpdateMetric(storager storages.MetricStorager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		Type := chi.URLParam(r, "type")
		id := chi.URLParam(r, "id")
		value := chi.URLParam(r, "value")

		metric, err := validateUpdate(Type, id, value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		if err := storager.Save(ctx, metric); err != nil {
			log.WithError(err).Error("failed to save metric")
			http.Error(w, "failed to save metric", http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// nolint: err113
func validateUpdate(mType, id, value string) (*schemas.Metric, error) {
	result := schemas.Metric{ID: id, MType: mType}

	if mType == schemas.Gauge {
		value, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, errors.New("invalid value")
		}

		result.Value = &value
	}

	if mType == schemas.Counter {
		value, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, errors.New("invalid delta")
		}

		result.Delta = &value
	}

	if err := result.Validate(); err != nil {
		return nil, err
	}

	return &result, nil
}
