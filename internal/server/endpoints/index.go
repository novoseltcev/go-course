package endpoints

import (
	"html/template"
	"net/http"

	"github.com/novoseltcev/go-course/internal/model"
)


func Index(counterStorage *MetricStorager[model.Counter], gaugeStorage *MetricStorager[model.Gauge]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			CounterMetrics []model.Metric[model.Counter]
			GaugeMetrics []model.Metric[model.Gauge]
		}{
			CounterMetrics: (*counterStorage).GetAll(),
			GaugeMetrics: (*gaugeStorage).GetAll(),
		}

		tmpl, _ := template.ParseFiles("templates/index.html")
        tmpl.Execute(w, data)
	}
}
