package endpoints

import (
	"html/template"
	"net/http"

	"github.com/novoseltcev/go-course/internal/types"
)


func Index(counterStorage *MetricStorager[types.Counter], gaugeStorage *MetricStorager[types.Gauge]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			CounterMetrics []types.Metric[types.Counter]
			GaugeMetrics []types.Metric[types.Gauge]
		}{
			CounterMetrics: (*counterStorage).GetAll(),
			GaugeMetrics: (*gaugeStorage).GetAll(),
		}

		tmpl, _ := template.ParseFiles("templates/index.html")
        tmpl.Execute(w, data)
	}
}
