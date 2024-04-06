package endpoints

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/novoseltcev/go-course/internal/model"
)


func GetOneMetric(counterStorage *MetricStorager[model.Counter], gaugeStorage *MetricStorager[model.Gauge]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")

		metricValue := ""
		switch metricType {
		case "gauge":
			result := (*gaugeStorage).GetByName(metricName)
			metricValue = strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", float64(result)), "0"), ".")
		case "counter":
			result := (*counterStorage).GetByName(metricName)
			metricValue = strconv.Itoa(int(result))
		default:
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		
		if metricValue == "" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		io.WriteString(w, metricValue)
	}
}
