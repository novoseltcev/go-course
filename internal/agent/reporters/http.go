package reporters

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/pkg/retry"
)

//go:generate mockgen -source=http.go -destination=./http_mock_test.go -package=reporters_test -typed
type (
	compressor interface {
		Compress(data []byte) ([]byte, error)
	}

	encryptor interface {
		Encrypt(data []byte) ([]byte, error)
	}

	hasher interface {
		GetHash(data []byte) ([]byte, error)
	}
)

type Option func(*HTTPReporter)

type HTTPReporter struct {
	c         *http.Client
	retryOpts *retry.Options
	baseURL   string
	enc       encryptor
	cmp       compressor
	hr        hasher
}

// WithRetry enables retrying of the request.
func WithRetry(opts retry.Options) Option {
	return func(r *HTTPReporter) {
		r.retryOpts = &opts
	}
}

// WithEncryption enables encryption of the request body.
func WithEncryption(enc encryptor) Option {
	return func(r *HTTPReporter) {
		r.enc = enc
	}
}

// WithCompression enables compression of the request body.
func WithCompression(cmp compressor) Option {
	return func(r *HTTPReporter) {
		r.cmp = cmp
	}
}

// WithCheckSum enables checksum signing of the request body.
func WithCheckSum(hr hasher) Option {
	return func(r *HTTPReporter) {
		r.hr = hr
	}
}

// NewHTTPClient creates a new ReportClient.
func NewHTTPReporter(c *http.Client, baseURL string, opts ...Option) *HTTPReporter {
	reportClient := &HTTPReporter{c: c, baseURL: baseURL}

	for _, opt := range opts {
		opt(reportClient)
	}

	return reportClient
}

func (rc *HTTPReporter) Report(ctx context.Context, metrics []schemas.Metric) error {
	if len(metrics) == 0 {
		return nil
	}

	b, err := rc.prepareData(metrics)
	if err != nil {
		return err
	}

	return rc.send(ctx, rc.baseURL+"/updates/", b)
}

func (rc *HTTPReporter) prepareData(metrics []schemas.Metric) ([]byte, error) {
	if len(metrics) == 0 {
		return nil, nil
	}

	b, err := json.Marshal(metrics)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal metrics: %w", err)
	}

	if rc.enc != nil {
		b, err = rc.enc.Encrypt(b)
		if err != nil {
			return nil, fmt.Errorf("cannot encrypt metrics: %w", err)
		}
	}

	if rc.cmp != nil {
		b, err = rc.cmp.Compress(b)
		if err != nil {
			return nil, fmt.Errorf("cannot compress metrics: %w", err)
		}
	}

	return b, nil
}

func (rc *HTTPReporter) send(ctx context.Context, url string, body []byte) error {
	if len(body) == 0 {
		return nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("X-Real-IP", getIP())

	if rc.cmp != nil {
		req.Header.Set("Content-Encoding", "gzip")
	}

	if rc.hr != nil {
		checkSum, hashErr := rc.hr.GetHash(body)
		if hashErr != nil {
			return fmt.Errorf("cannot get hash: %w", hashErr)
		}

		req.Header.Set("Hashsha256", hex.EncodeToString(checkSum))
	}

	resp, err := retry.DoWithData(ctx, func() (*http.Response, error) {
		return rc.c.Do(req)
	}, rc.retryOpts)
	if err != nil {
		return fmt.Errorf("error during send request: %w", err)
	}

	defer resp.Body.Close()

	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		return fmt.Errorf("error during read request body: %w", err)
	}

	log.WithFields(log.Fields{"url": req.URL, "status": resp.Status}).Info("report successfully sent")

	return nil
}

func getIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}
