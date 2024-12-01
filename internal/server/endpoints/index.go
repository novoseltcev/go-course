package endpoints

import (
	"embed"
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/storages"
)

//go:embed "templates/*"
var templatesFS embed.FS

func Index(storager storages.MetricStorager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics, err := storager.GetAll(r.Context())
		if err != nil {
			log.WithError(err).Error("failed to get metrics")
			http.Error(w, "failed to get metrics", http.StatusInternalServerError)

			return
		}

		tmpl := template.Must(template.ParseFS(templatesFS, "templates/index.gohtml"))

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		if err := tmpl.Execute(w, metrics); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
