package reporters

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/pkg/retry"
)

type Requester interface {
	Do(r *http.Request) (*http.Response, error)
}

type APIReporter struct {
	client  Requester
	baseURL string
	hashKey string
}

func NewAPIReporter(client Requester, baseURL, hashKey string) *APIReporter {
	return &APIReporter{
		client:  client,
		baseURL: baseURL,
		hashKey: hashKey,
	}
}

func (r *APIReporter) Report(ctx context.Context, metrics []schemas.Metric) error {
	buf, err := compress(metrics)
	if err != nil {
		log.WithError(err).Error("cannot compress metrics")

		return err
	}

	checkSum := r.calculateCheckSum(buf.Bytes())
	url := r.baseURL + "/updates/"

	attempts := []time.Duration{time.Second, 3 * time.Second, 5 * time.Second}

	response, err := retry.DoWithData(ctx, func() (*http.Response, error) {
		return r.send(ctx, checkSum, url, buf)
	}, &retry.Options{
		Retries:  uint(len(attempts)),
		Attempts: attempts,
	})
	if err != nil {
		log.WithError(err).WithField("url", url).Error("Error during send request")

		return err
	}

	if response.Body != nil {
		defer response.Body.Close()

		if _, err := io.Copy(io.Discard, response.Body); err != nil {
			log.WithError(err).WithField("url", url).Error("Error during read request body")

			return err
		}
	}

	log.WithFields(log.Fields{"url": url, "status": response.Status}).Info("request sent")

	return nil
}

func (r *APIReporter) calculateCheckSum(data []byte) []byte {
	if r.hashKey == "" {
		return nil
	}

	h := hmac.New(sha256.New, []byte(r.hashKey))
	h.Write(data)

	return h.Sum(nil)
}

func (r *APIReporter) send(ctx context.Context, checkSum []byte, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	if checkSum != nil {
		req.Header.Set("Hashsha256", hex.EncodeToString(checkSum))
	}

	return r.client.Do(req)
}

func compress(metrics []schemas.Metric) (*bytes.Buffer, error) {
	result, err := json.Marshal(metrics)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	gzw := gzip.NewWriter(buf)

	if _, err := gzw.Write(result); err != nil {
		return nil, err
	}

	gzw.Close()

	return buf, nil
}
