package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"net/http"
)

type hashWriter struct {
	http.ResponseWriter
	h hash.Hash
}

func (w hashWriter) Write(p []byte) (int, error) {
	if written, err := w.h.Write(p); err != nil {
		return written, err
	}

	w.Header().Set("Hashsha256", hex.EncodeToString(w.h.Sum(nil)))

	return w.ResponseWriter.Write(p)
}

func Sign(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wh := hmac.New(sha256.New, []byte(key))
			next.ServeHTTP(&hashWriter{w, wh}, r)
		})
	}
}
