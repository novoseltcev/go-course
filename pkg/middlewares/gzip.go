package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type compressReader struct {
	body         io.ReadCloser
	decompressor *gzip.Reader
}

func (c compressReader) Read(p []byte) (int, error) {
	return c.decompressor.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.body.Close(); err != nil {
		return err
	}

	return c.decompressor.Close()
}

type compressWriter struct {
	http.ResponseWriter
	compressor *gzip.Writer
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.compressor.Write(p)
}

// Gzip is compression middleware.
//
// It will compress the response if the client accepts gzip encoding.
// It will decompress the request if the client sends gzip encoding.
func Gzip(next http.Handler) http.Handler {
	wrapper := func(w http.ResponseWriter, r *http.Request) {
		contentEncoding := r.Header.Get("Content-Encoding")
		if isCompressed := strings.Contains(contentEncoding, "gzip"); isCompressed {
			decompressor, err := gzip.NewReader(r.Body)
			if err != nil {
				log.WithError(err).Error("cannot decompress body")
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

				return
			}

			cr := &compressReader{
				body:         r.Body,
				decompressor: decompressor,
			}
			defer cr.Close()
			r.Body = cr
		}

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)

			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		gw := gzip.NewWriter(w)

		defer gw.Close()

		next.ServeHTTP(&compressWriter{ResponseWriter: w, compressor: gw}, r)
	}

	return http.HandlerFunc(wrapper)
}
