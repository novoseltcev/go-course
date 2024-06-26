package workers

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"time"

	json "github.com/mailru/easyjson"
	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/model"
	"github.com/novoseltcev/go-course/internal/schema"
)


type Client interface {
	Do(*http.Request) (*http.Response, error)
}

func SendMetrics(jobsCh <-chan model.Metric, rateLimit time.Duration, client Client, baseURL, hashKey string) {
	logger := log.WithField("workerName", "SendMetrics")
	logger.Info("start worker")

	lastSent := time.Now()
	body := make([]schema.Metrics, 0)

	for metric := range jobsCh {
		body = append(body, schema.Metrics{ID: metric.Name, MType: metric.Type, Value: metric.Value, Delta: metric.Delta})

		if time.Since(lastSent) > rateLimit {
			if err := send(client, baseURL, hashKey, body); err != nil {
				logger.WithError(err).Error("crash worker")
			} else {
				body = make([]schema.Metrics, 0)
			}
			lastSent = time.Now()
		}
	}
}

func send(c Client, baseURL, hashKey string, metrics schema.MetricsSlice) error {
	result, err := json.Marshal(metrics); 
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	if _, err := zb.Write(result); err != nil {
		return err
	}
	zb.Close()

	var checkSum []byte;
	if hashKey != "" {
		h := hmac.New(sha256.New, []byte(hashKey))
		h.Write(buf.Bytes())
		checkSum = h.Sum(nil)
	}

	url := baseURL + "/updates/"
	response, err := post(c, url, checkSum, buf)
	if err != nil {
		fmt.Printf("Error during send request to %s\n", url)
		return err
	}
	
	if response.Body != nil {
		defer response.Body.Close()
		if _, err := io.Copy(io.Discard, response.Body); err != nil {
			fmt.Printf("Error during read request body from response to %s\n", url)
			return err
		}
	}

	fmt.Printf("Sended request to %s with code %s\n", url, response.Status)
	return nil
}

func post(c Client, url string, checkSum []byte, body io.Reader) (*http.Response, error) {
	var (
        err error
        response *http.Response
    )
	timeouts := []time.Duration{time.Second, 3 * time.Second, 5 * time.Second}
	retries := len(timeouts)
    for retries > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), timeouts[len(timeouts) - retries])
		defer cancel()

        req, reqErr := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
		if reqErr != nil {
			return nil, reqErr
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Encoding", "gzip")
		if checkSum != nil {
			req.Header.Set("HashSHA256", fmt.Sprintf("%x", checkSum))
		}

		response, err = c.Do(req)
        if err != nil {
			fmt.Println("retry send metrics")
            retries -= 1
        } else {
            break
        }
    }

	return response, err
}
