package endpoints

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/novoseltcev/go-course/internal/model"
	"github.com/novoseltcev/go-course/internal/server/storage"
	"github.com/novoseltcev/go-course/internal/utils"
)


func Index(storage storage.MetricStorager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		
		metrics, err := utils.RetryPgSelect(ctx, func() ([]model.Metric, error) {
			return storage.GetAll(ctx)
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Printf("%d", len(metrics))

		tmpl, _ := template.ParseFiles("templates/index.html")
        w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, metrics)
		
	}
}
