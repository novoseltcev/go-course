package endpoints

import (
	"html/template"
	"net/http"
	"github.com/novoseltcev/go-course/internal/server/storage"
)


func Index(counterStorage *storage.Storage[storage.Counter], gaugeStorage *storage.Storage[storage.Gauge]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			CounterMetrics []storage.Metric[storage.Counter]
			GaugeMetrics []storage.Metric[storage.Gauge]
		}{
			CounterMetrics: (*counterStorage).GetAll(),
			GaugeMetrics: (*gaugeStorage).GetAll(),
		}

		tmpl, _ := template.ParseFiles("templates/index.html")
        tmpl.Execute(w, data)
	}
}
