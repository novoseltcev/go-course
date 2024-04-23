package endpoints

import (
	"html/template"
	"net/http"

	"github.com/novoseltcev/go-course/internal/model"
	"github.com/novoseltcev/go-course/internal/server/storage"
)


func Index(counterStorage *storage.MetricStorager[model.Counter], gaugeStorage *storage.MetricStorager[model.Gauge]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
	
		data := struct {
			CounterMetrics []model.Metric[model.Counter]
			GaugeMetrics []model.Metric[model.Gauge]
		}{
			CounterMetrics: (*counterStorage).GetAll(ctx),
			GaugeMetrics: (*gaugeStorage).GetAll(ctx),
		}

		tmpl, _ := template.ParseFiles("templates/index.html")
        w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, data)
		
	}
}
