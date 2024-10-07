package endpoints

import (
	"html/template"
	"net/http"

	"github.com/novoseltcev/go-course/internal/server/services"
	"github.com/novoseltcev/go-course/internal/server/storage"
)

func Index(storage storage.MetricStorager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		metrics, err := services.GetAllMetric(ctx, storage, pgRetries)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		if err := tmpl.Execute(w, metrics); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
