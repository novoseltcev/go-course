package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type compressReader struct {
    body io.ReadCloser
    decompressor *gzip.Reader
}

func (c compressReader) Read(p []byte) (n int, err error) {
    return c.decompressor.Read(p)
}

func (c *compressReader) Close() error {
    if err := c.body.Close(); err != nil {
        return err
    }
    return c.decompressor.Close()
} 


type compressWriter struct {
    writer http.ResponseWriter
    compressor *gzip.Writer
}

func (c *compressWriter) Write(p []byte) (int, error) {
    return c.compressor.Write(p)
}

func (c *compressWriter) Header() http.Header {
    return c.writer.Header()
}

func (c *compressWriter) WriteHeader(statusCode int) {
    if statusCode < 300 {
        c.writer.Header().Set("Content-Encoding", "gzip")
    }
    c.writer.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
    return c.compressor.Close()
}

func Gzip(handler http.Handler) http.Handler {
    wrapper := func(w http.ResponseWriter, r *http.Request) {        
        
        contentEncoding := r.Header.Get("Content-Encoding")
        if isCompressed := strings.Contains(contentEncoding, "gzip"); isCompressed {
            decompressor, err := gzip.NewReader(r.Body)
			if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
			}
			cr := compressReader{
				body: r.Body,
				decompressor: decompressor,
			}
            r.Body = &cr
			defer cr.Close()
        }

        ow := w
        acceptEncoding := r.Header.Get("Accept-Encoding")
        contentType := r.Header.Get("Content-Type")
        supportGzip := strings.Contains(acceptEncoding, "gzip")
        compessableType := strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/html")
        
        if supportGzip && compessableType {
            cw := &compressWriter{
                writer:  w,
                compressor: gzip.NewWriter(w),
            }
            ow = cw
            defer cw.Close()
        }

        handler.ServeHTTP(ow, r)
    }

	return http.HandlerFunc(wrapper)
}
