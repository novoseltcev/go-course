package endpoints

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/internal/storages"
)

func GetOneMetric(storager storages.MetricStorager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		reqData := schemas.MetricIdentifier{
			ID:    chi.URLParam(r, "id"),
			MType: chi.URLParam(r, "type"),
		}

		if err := reqData.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		metric, err := storager.GetOne(ctx, reqData.ID, reqData.MType)
		if err != nil {
			if errors.Is(err, storages.ErrNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				log.WithField("id", reqData.ID).WithField("type", reqData.MType).WithError(err).Error("failed to get metric")
				http.Error(w, "failed to get metric", http.StatusInternalServerError)
			}

			return
		}

		result, err := serialize(metric)
		if err != nil {
			log.WithError(err).Error("failed to serialize metric")
			http.Error(w, "failed to serialize", http.StatusInternalServerError)

			return
		}

		if _, err := io.WriteString(w, result); err != nil {
			log.WithError(err).Error("failed to write metric")
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func serialize(metric *schemas.Metric) (string, error) {
	if err := metric.Validate(); err != nil {
		return "", err
	}

	switch metric.MType {
	case schemas.Gauge:
		return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", *metric.Value), "0"), "."), nil
	case schemas.Counter:
		return strconv.Itoa(int(*metric.Delta)), nil
	default:
		return "", schemas.ErrInvalidType
	}
}
