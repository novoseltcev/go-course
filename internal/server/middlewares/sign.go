package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"hash"
	"net/http"
)


type hashWriter struct {
    http.ResponseWriter
    h hash.Hash
}

func (w hashWriter) Write(p []byte) (int, error) {    
    if writen, err := w.h.Write(p); err != nil {
        return writen, err
    }

    return w.ResponseWriter.Write(p)
}

func Sign(key string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        wrapper := func(w http.ResponseWriter, r *http.Request) {
            wh := hmac.New(sha256.New, []byte(key))
            next.ServeHTTP(&hashWriter{w, wh}, r)

            w.Header().Set("HashSHA256", fmt.Sprintf("%x", wh.Sum(nil)))
        }

        return http.HandlerFunc(wrapper)
    }
}
