package endpoints

import (
	"net/http"

	"github.com/novoseltcev/go-course/internal/storages"
)

func Ping(storager storages.MetricStorager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := storager.Ping(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
