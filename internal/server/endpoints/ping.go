package endpoints

import (
	"net/http"

	"github.com/novoseltcev/go-course/internal/server/storage"
)

func Ping(storage storage.MetricStorager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := storage.Ping(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
