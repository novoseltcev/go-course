package endpoints

import (
	"html/template"
	"net/http"

	"github.com/novoseltcev/go-course/internal/model"
	"github.com/novoseltcev/go-course/internal/server/storage"
)


func Index(counterStorage *storage.MetricStorager[model.Counter], gaugeStorage *storage.MetricStorager[model.Gauge]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			CounterMetrics []model.Metric[model.Counter]
			GaugeMetrics []model.Metric[model.Gauge]
		}{
			CounterMetrics: (*counterStorage).GetAll(),
			GaugeMetrics: (*gaugeStorage).GetAll(),
		}

		tmpl, _ := template.ParseFiles("templates/index.html")
        w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, data)
		
	}
}
