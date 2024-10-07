package middlewares

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var ErrInvalidCheckSum = errors.New("invalid check sum")

func CheckSum(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		wrapper := func(w http.ResponseWriter, r *http.Request) {
			requestCheckSum, err := hex.DecodeString(r.Header.Get("Hashsha256"))
			if err != nil {
				log.WithError(err).Errorf("HashSHA256 not decodable to string")
				http.Error(w, err.Error(), http.StatusBadRequest)

				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.WithError(err).Error("cannot read body")
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

				return
			}

			r.Body.Close()

			r.Body = io.NopCloser(bytes.NewBuffer(body))

			h := hmac.New(sha256.New, []byte(key))
			h.Write(body)
			computedCheckSum := h.Sum(nil)

			if !hmac.Equal(requestCheckSum, computedCheckSum) {
				log.WithFields(
					log.Fields{
						"got":  requestCheckSum,
						"want": computedCheckSum,
					},
				).Error(ErrInvalidCheckSum.Error())
				http.Error(w, ErrInvalidCheckSum.Error(), http.StatusBadRequest)

				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(wrapper)
	}
}
