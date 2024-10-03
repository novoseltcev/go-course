
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
	ch := make(chan model.Metric)
	var client Client = ClientMock{}
	baseURL := "http://0.0.0.0:8080"

	go SendMetrics(ch, 1, client, baseURL, "secret-key")

	value := 123.321
	var delta int64 = 2
	ch <- model.Metric{Type: "gauge", Name: "Some", Value: &value}
	ch <- model.Metric{Type: "counter", Name: "Some", Delta: &delta}
	close(ch)
}

func BenchmarkSendMetrics(b *testing.B) {
	ch := make(chan model.Metric)
	client := ClientMock{}
	baseURL := "http://0.0.0.0:8080"

	go SendMetrics(ch, 1, client, baseURL, "secret-key")

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		value := 123.321
		var delta int64 = 2
		ch <- model.Metric{Type: "gauge", Name: "Some", Value: &value}
		ch <- model.Metric{Type: "counter", Name: "Some", Delta: &delta}
	}
}
