package middlewares

import (
	"bytes"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

//go:generate mockgen -source=decrypt.go -destination=./mock_test.go -package=middlewares_test
type decryptor interface {
	Decrypt(b []byte) ([]byte, error)
}

func Decrypt(dec decryptor) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Body != nil {
				b, err := io.ReadAll(r.Body)
				if err != nil {
					log.WithError(err).Error("failed to read body")
					http.Error(w, "failed to read body", http.StatusInternalServerError)

					return
				}

				decrypted, err := dec.Decrypt(b)
				if err != nil {
					log.WithError(err).Error("failed to decrypt")
					http.Error(w, "failed to decrypt", http.StatusInternalServerError)

					return
				}

				r.Body = io.NopCloser(bytes.NewBuffer(decrypted))
			}

			next.ServeHTTP(w, r)
		})
	}
}
