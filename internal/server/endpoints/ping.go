package endpoints

import (
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/novoseltcev/go-course/internal/utils"
)

func Ping(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		ctx := r.Context()
		
		err := utils.RetryPgExec(ctx, func() error {
			return db.PingContext(ctx)
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
