package workers

import (
	"net/http"
	"testing"

	"github.com/novoseltcev/go-course/internal/model"
)

type ClientMock struct {}

func (c ClientMock) Do (*http.Request) (*http.Response, error) {
    return &http.Response{}, nil
}

func TestSendMetrics(t *testing.T) {
	var client Client = ClientMock{}
	baseURL := "http://0.0.0.0:8080"

	ch := make(chan model.Metric)

	go SendMetrics(ch, 1, client, baseURL, "secret-key")

	var value float64 = 123.321
	var delta int64 = 2
	ch <- model.Metric{Type: "gauge", Name: "Some", Value: &value}
	ch <- model.Metric{Type: "counter", Name: "Some", Delta: &delta}
	close(ch)
}
