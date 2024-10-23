package endpoints

import (
	"context"
	"net/http"
)

type Pinger interface {
	Ping(ctx context.Context) error
}

func Ping(p Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := p.Ping(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
