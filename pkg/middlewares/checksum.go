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

// CheckSum is a middleware that checks the checksum of the request body.
//
// The middleware expects the key to be passed as a parameter.
// It expects the checksum to be in the header "Hashsha256".
// The middleware will check if the checksum is equal to the computed checksum.
func CheckSum(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		wrapper := func(w http.ResponseWriter, r *http.Request) {
			if r.Body != nil {
				requestCheckSum, err := hex.DecodeString(r.Header.Get("Hashsha256"))
				if err != nil {
					log.WithError(err).Errorf("Hashsha256 not decodable to string")
					http.Error(w, err.Error(), http.StatusBadRequest)

					return
				}

				body, err := copyBody(r)
				if err != nil {
					log.WithError(err).Error("cannot read body")
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

					return
				}

				h := hmac.New(sha256.New, []byte(key))
				h.Write(body)
				computedCheckSum := h.Sum(nil)

				if !hmac.Equal(requestCheckSum, computedCheckSum) {
					log.WithField("got", requestCheckSum).WithField("want", computedCheckSum).Error(ErrInvalidCheckSum.Error())
					http.Error(w, ErrInvalidCheckSum.Error(), http.StatusBadRequest)

					return
				}
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(wrapper)
	}
}

func copyBody(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	r.Body.Close()
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	return body, nil
}
