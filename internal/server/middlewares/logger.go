package middlewares

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type (
	responseData struct {
		status int
		size int
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

func Logger(handler http.Handler) http.Handler {
	wrapper := func(w http.ResponseWriter, r *http.Request) {
		lw := loggingResponseWriter {
			ResponseWriter: w,
			responseData: responseData {
				status: 0,
				size: 0,
			},
		}
		start := time.Now()
		handler.ServeHTTP(&lw, r)
		elapsed := time.Since(start)
		log.Infof("%s %s - %d %dB in %.3fµs", r.Method, r.URL, lw.responseData.status, lw.responseData.size, float64(elapsed) / 1000)
	}

	return http.HandlerFunc(wrapper)
}
