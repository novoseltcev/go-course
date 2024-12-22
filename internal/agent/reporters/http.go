package reporters

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/pkg/retry"
)

//go:generate mockgen -source=http.go -destination=./http_mock_test.go -package=reporters_test
type (
	compressor interface {
		Compress(data []byte) ([]byte, error)
	}

	hasher interface {
		GetHash(data []byte) ([]byte, error)
	}
)

type Option func(*ReportClient)

type ReportClient struct {
	c         *http.Client
	retryOpts *retry.Options
	baseURL   string
	cmp       compressor
	hr        hasher
}

// WithRetry enables retrying of the request.
func WithRetry(opts retry.Options) Option {
	return func(r *ReportClient) {
		r.retryOpts = &opts
	}
}

// WithCompression enables compression of the request body.
func WithCompression(cmp compressor) Option {
	return func(r *ReportClient) {
		r.cmp = cmp
	}
}

// WithCheckSum enables checksum signing of the request body.
func WithCheckSum(hr hasher) Option {
	return func(r *ReportClient) {
		r.hr = hr
	}
}

// NewHTTPClient creates a new ReportClient.
func NewHTTPClient(c *http.Client, baseURL string, opts ...Option) *ReportClient {
	reportClient := &ReportClient{c: c, baseURL: baseURL} //nolint:exhaustruct

	for _, opt := range opts {
		opt(reportClient)
	}

	return reportClient
}

func (rc *ReportClient) Report(ctx context.Context, metrics []schemas.Metric) error {
	b, err := rc.prepareData(metrics)
	if err != nil {
		return err
	}

	return rc.send(ctx, rc.baseURL+"/updates/", b)
}

func (rc *ReportClient) prepareData(metrics []schemas.Metric) ([]byte, error) {
	b, err := json.Marshal(metrics)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal metrics: %w", err)
	}

	if rc.cmp != nil {
		b, err = rc.cmp.Compress(b)
		if err != nil {
			return nil, fmt.Errorf("cannot compress metrics: %w", err)
		}
	}

	return b, nil
}

func (rc *ReportClient) send(ctx context.Context, url string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

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
