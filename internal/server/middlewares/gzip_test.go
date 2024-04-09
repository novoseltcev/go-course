package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)


func webhook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
    b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(b)
	w.Header().Set("Content-Type", "application/json")
} 

func TestGzipCompression(t *testing.T) {
    handler := http.HandlerFunc(webhook)
    srv := httptest.NewServer(Gzip(handler))
    defer srv.Close()
    
    body := `{"ping": "pong"}`

	t.Run("without_gzip", func(t *testing.T) {
		buf := bytes.NewBufferString(body)
        r := httptest.NewRequest("POST", srv.URL, buf)
        r.RequestURI = ""
        r.Header.Set("Accept-Encoding", "")
        
        resp, err := http.DefaultClient.Do(r)
        require.NoError(t, err)
        require.Equal(t, http.StatusOK, resp.StatusCode)
        
        defer resp.Body.Close()
        
        b, err := io.ReadAll(resp.Body)
        require.NoError(t, err)
        require.JSONEq(t, body, string(b))
    })

    t.Run("sends_gzip", func(t *testing.T) {
        buf := bytes.NewBuffer(nil)
        zb := gzip.NewWriter(buf)
        _, err := zb.Write([]byte(body))
        require.NoError(t, err)
        err = zb.Close()

        require.NoError(t, err)
        
        r := httptest.NewRequest("POST", srv.URL, buf)
        r.RequestURI = ""
        r.Header.Set("Content-Encoding", "gzip")
        r.Header.Set("Accept-Encoding", "")
        
        resp, err := http.DefaultClient.Do(r)
        require.NoError(t, err)
        require.Equal(t, http.StatusOK, resp.StatusCode)
        
        defer resp.Body.Close()
        
        b, err := io.ReadAll(resp.Body)
        require.NoError(t, err)
        require.JSONEq(t, body, string(b))
    })

    t.Run("accepts_gzip", func(t *testing.T) {
        buf := bytes.NewBufferString(body)
        r := httptest.NewRequest("POST", srv.URL, buf)
        r.RequestURI = ""

        r.Header.Set("Content-Type", "application/json")
        r.Header.Set("Accept-Encoding", "gzip")
        
        resp, err := http.DefaultClient.Do(r)
        require.NoError(t, err)
        require.Equal(t, http.StatusOK, resp.StatusCode)
        
        defer resp.Body.Close()
        
        zr, err := gzip.NewReader(resp.Body)
        require.NoError(t, err)
        
        b, err := io.ReadAll(zr)
        require.NoError(t, err)
        require.JSONEq(t, body, string(b))
    })
}
