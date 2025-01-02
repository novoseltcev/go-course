package middlewares

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size

	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// Logger middleware.
//
// It will log the method url, status code and the size of the response body with the elapsed time in nanoseconds.
func Logger(handler http.Handler) http.Handler {
	wrapper := func(w http.ResponseWriter, r *http.Request) {
		lrw := loggingResponseWriter{ResponseWriter: w, responseData: responseData{}}
		start := time.Now()

		handler.ServeHTTP(&lrw, r)

		elapsed := time.Since(start)
		log.Infof(
			"%s %s - %d %dB in %.3fµs",
			r.Method,
			r.URL,
			lrw.responseData.status,
			lrw.responseData.size,
			float64(elapsed)/float64(time.Microsecond),
		)
	}

	return http.HandlerFunc(wrapper)
}
