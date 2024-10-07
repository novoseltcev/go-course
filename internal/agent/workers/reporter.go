package workers

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"time"

	json "github.com/mailru/easyjson"
	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/schema"
)

type Client interface {
	Do(r *http.Request) (*http.Response, error)
}

func SendMetrics(
	ctx context.Context,
	jobsCh <-chan schema.Metric,
	c Client,
	rateLimit time.Duration,
	baseURL, hashKey string,
) {
	logger := log.WithField("workerName", "SendMetrics")
	logger.Info("start worker")

	lastSent := time.Now()
	body := make([]schema.Metric, 0)

	for metric := range jobsCh {
		select {
		case <-ctx.Done():
			return
		default:
		}

		body = append(body, schema.Metric{ID: metric.ID, MType: metric.MType, Value: metric.Value, Delta: metric.Delta})

		if time.Since(lastSent) > rateLimit {
			if err := send(ctx, c, baseURL, hashKey, body); err != nil {
				logger.WithError(err).Error("crash worker")
			} else {
				body = make([]schema.Metric, 0)
			}

			lastSent = time.Now()
		}
	}
}

func send(
	ctx context.Context,
	c Client,
	baseURL, hashKey string,
	metrics schema.MetricSlice,
) error {
	result, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	gzw := gzip.NewWriter(buf)

	if _, err := gzw.Write(result); err != nil {
		return err
	}

	gzw.Close()

	var checkSum []byte

	if hashKey != "" {
		h := hmac.New(sha256.New, []byte(hashKey))
		h.Write(buf.Bytes())
		checkSum = h.Sum(nil)
	}

	url := baseURL + "/updates/"

	response, err := post(ctx, c, url, checkSum, buf)
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

func post(ctx context.Context, c Client, url string, checkSum []byte, body io.Reader) (*http.Response, error) {
	var (
		err      error
		response *http.Response
	)

	timeouts := []time.Duration{time.Second, 3 * time.Second, 5 * time.Second}
	retries := len(timeouts)

	for retries > 0 {
		ctx, cancel := context.WithTimeout(ctx, timeouts[len(timeouts)-retries])
		defer cancel()

		req, reqErr := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
		if reqErr != nil {
			return nil, reqErr
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Encoding", "gzip")

		if checkSum != nil {
			req.Header.Set("Hashsha256", hex.EncodeToString(checkSum))
		}

		response, err = c.Do(req)
		if err != nil {
			log.WithError(err).Error("retry send metrics")

			retries--
		} else {
			break
		}
	}

	return response, err
}
